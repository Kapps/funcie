package bastion

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports"
	"github.com/google/uuid"
	"log/slog"
	"syscall"
)

type handler struct {
	registry       funcie.ApplicationRegistry
	appClient      ApplicationClient
	consumer       funcie.Consumer
	hostTranslator HostTranslator
}

// NewHandler creates a new Handler that can register and unregister applications and forward requests.
func NewHandler(
	registry funcie.ApplicationRegistry,
	appClient ApplicationClient,
	consumer funcie.Consumer,
	hostTranslator HostTranslator,
) transports.MessageHandler {
	return &handler{
		registry:       registry,
		appClient:      appClient,
		consumer:       consumer,
		hostTranslator: hostTranslator,
	}
}

func (h *handler) Register(ctx context.Context, message messages.RegistrationMessage) (*messages.RegistrationResponse, error) {
	application := funcie.NewApplication(message.Payload.Name, message.Payload.Endpoint)
	translatedHost, err := h.hostTranslator.TranslateLocalHostToResolvedHost(ctx, application.Endpoint.Host)
	if err != nil {
		return nil, fmt.Errorf("translate local host %v to resolved host: %w", application.Endpoint.Host, err)
	}

	application.Endpoint.Host = translatedHost

	err = h.registry.Register(ctx, application)
	if err != nil {
		return nil, fmt.Errorf("register application %v: %w", application, err)
	}

	if err := h.consumer.Subscribe(ctx, application.Name, h.onConsumerMessageReceived); err != nil {
		return nil, fmt.Errorf("subscribe to application %v: %w", application, err)
	}

	registrationId := uuid.New()
	slog.InfoContext(ctx, "registered application", "application", application, "registrationId", registrationId)

	responsePayload := messages.NewRegistrationResponsePayload(registrationId)
	return funcie.NewResponseWithPayload(message.ID, responsePayload, nil), nil
}

func (h *handler) Deregister(ctx context.Context, message messages.DeregistrationMessage) (*messages.DeregistrationResponse, error) {
	applicationName := message.Payload.Name
	err := h.registry.Unregister(ctx, applicationName)
	if err != nil {
		return nil, fmt.Errorf("unregister application %v: %w", applicationName, err)
	}

	if err := h.consumer.Unsubscribe(ctx, applicationName); err != nil {
		return nil, fmt.Errorf("unsubscribe from application %v: %w", applicationName, err)
	}

	responsePayload := messages.NewDeregistrationResponsePayload()
	return funcie.NewResponseWithPayload(message.ID, responsePayload, nil), nil
}

func (h *handler) ForwardRequest(ctx context.Context, request messages.ForwardRequestMessage) (*messages.ForwardRequestResponse, error) {
	app, err := h.registry.GetApplication(ctx, request.Application)
	if err == funcie.ErrApplicationNotFound {
		slog.WarnContext(ctx, "application not found in client registry", "application", request.Application)
		// TODO: What should we return here?
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("getting application %v: %w", request.Application, err)
	}

	marshaled, err := funcie.MarshalMessagePayload[messages.ForwardRequestMessage](request)
	resp, err := h.appClient.ProcessRequest(ctx, *app, marshaled)
	if err != nil {
		return nil, fmt.Errorf("process request %v: %w", request.ID, err)
	}

	unmarshaled, err := funcie.UnmarshalResponsePayload[messages.ForwardRequestResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response payload: %w", err)
	}

	return unmarshaled, nil
}

func (h *handler) onConsumerMessageReceived(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	if message.Kind != messages.MessageKindForwardRequest {
		slog.WarnContext(ctx, "ignoring invalid message kind", "kind", message.Kind)
		return nil, nil
	}

	// TODO: More or less a reimplentation of MessageProcessor -- needs some refactoring.

	app, err := h.registry.GetApplication(ctx, message.Application)
	if err == funcie.ErrApplicationNotFound {
		slog.WarnContext(ctx, "application not found in client registry", "application", message.Application)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting application %v: %w", message.Application, err)
	}

	resp, err := h.appClient.ProcessRequest(ctx, *app, message)
	if errors.Is(err, syscall.ECONNREFUSED) {
		slog.WarnContext(ctx, "application not available", "application", message.Application)
		return funcie.NewResponse(message.ID, nil, funcie.ErrNoActiveConsumer), nil
	}
	if err != nil {
		return nil, fmt.Errorf("forward request: %w", err)
	}

	marshaled, err := funcie.MarshalResponsePayload(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response payload: %w", err)
	}

	return marshaled, nil
}
