package bastion

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

// TODO: This was never actually needed. Currently unused, will be removed soon.

type forwardingApplicationRegistry struct {
	underlying funcie.ApplicationRegistry
	publisher  funcie.Publisher
}

// NewForwardingApplicationRegistry creates a new ApplicationRegistry that forwards requests to a server bastion.
func NewForwardingApplicationRegistry(underlying funcie.ApplicationRegistry, publisher funcie.Publisher) funcie.ApplicationRegistry {
	return &forwardingApplicationRegistry{
		underlying: underlying,
		publisher:  publisher,
	}
}

func (f *forwardingApplicationRegistry) Register(ctx context.Context, application *funcie.Application) error {
	err := f.underlying.Register(ctx, application)
	if err != nil {
		return fmt.Errorf("failed to register application: %w", err)
	}

	payload := messages.NewRegistrationRequestPayload(application.Name, application.Endpoint)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling payload: %w", err)
	}

	message := funcie.NewMessage(application.Name, messages.MessageKindRegister, payloadBytes)

	resp, err := f.publisher.Publish(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to publish application registered event: %w", err)
	}

	if err := resp.Error; err != nil {
		return fmt.Errorf("failed to register application %v: %w", application.Name, err)
	}

	return nil
}

func (f *forwardingApplicationRegistry) Unregister(ctx context.Context, applicationName string) error {
	err := f.underlying.Unregister(ctx, applicationName)
	if err != nil {
		return fmt.Errorf("failed to unregister application: %w", err)
	}

	payload := messages.NewDeregistrationRequestPayload(applicationName)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling payload: %w", err)
	}

	message := funcie.NewMessage(applicationName, messages.MessageKindDeregister, payloadBytes)

	resp, err := f.publisher.Publish(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to publish application unregistered event: %w", err)
	}

	if err := resp.Error; err != nil {
		return fmt.Errorf("failed to unregister application %v: %w", applicationName, err)
	}

	return nil
}

func (f *forwardingApplicationRegistry) GetApplication(ctx context.Context, applicationName string) (*funcie.Application, error) {
	app, err := f.underlying.GetApplication(ctx, applicationName)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return app, nil
}
