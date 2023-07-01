package bastion

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"golang.org/x/exp/slog"
)

// Handler allows the handling of incoming valid Bastion requests.
type Handler interface {
	// Register registers the given application.
	Register(ctx context.Context, application *funcie.Application) error
	// Unregister unregisters the application with the given name.
	Unregister(ctx context.Context, applicationName string) error
	// ForwardRequest forwards the given request to the application specified in the request.
	ForwardRequest(ctx context.Context, request *messages.Message) (*funcie.Response, error)
}

type handler struct {
	registry  funcie.ApplicationRegistry
	publisher funcie.Publisher
}

// NewHandler creates a new Handler that can register and unregister applications and forward requests.
func NewHandler(registry funcie.ApplicationRegistry, publisher funcie.Publisher) Handler {
	return &handler{
		registry: registry,
	}
}

func (h *handler) Register(ctx context.Context, application *funcie.Application) error {
	err := h.registry.Register(ctx, application)
	if err != nil {
		return fmt.Errorf("register application %v: %w", application, err)
	}

	return nil
}

func (h *handler) Unregister(ctx context.Context, applicationName string) error {
	err := h.registry.Unregister(ctx, applicationName)
	if err != nil {
		return fmt.Errorf("unregister application %v: %w", applicationName, err)
	}

	return nil
}

func (h *handler) ForwardRequest(ctx context.Context, request *messages.Message) (*funcie.Response, error) {
	app, err := h.registry.GetApplication(ctx, request.Application)
	if err == funcie.ErrApplicationNotFound {
		slog.WarnCtx(ctx, "application not found in client registry", "application", request.Application)
		// TODO: What should we return here?
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("getting application %v: %w", request.Application, err)
	}

	// Server is forwarding us a request -- forward it as a simple http call to the process.

	return resp, nil
}
