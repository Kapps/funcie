package bastion

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

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

	message := messages.NewMessage(application.Name, funcie.MessageKindRegistration, Register)
	err = f.publisher.Publish(ctx, funcie.NewApplicationRegisteredEvent(application))
	if err != nil {
		return fmt.Errorf("failed to publish application registered event: %w", err)
	}
}

func (f *forwardingApplicationRegistry) Unregister(ctx context.Context, applicationName string) error {
	//TODO implement me
	panic("implement me")
}

func (f *forwardingApplicationRegistry) GetApplication(ctx context.Context, applicationName string) (*funcie.Application, error) {
	//TODO implement me
	panic("implement me")
}
