package transports

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie/messages"
)

// MessageHandler allows the handling of incoming valid Bastion requests.
type MessageHandler interface {
	// Register registers the given application.
	Register(ctx context.Context, message messages.RegistrationMessage) (*messages.RegistrationResponse, error)
	// Deregister removes the registration of the application with the given name.
	Deregister(ctx context.Context, message messages.DeregistrationMessage) (*messages.DeregistrationResponse, error)
	// ForwardRequest forwards the given request to the application specified in the request.
	ForwardRequest(ctx context.Context, message messages.ForwardRequestMessage) (*messages.ForwardRequestResponse, error)
}

// Implementation done on both sides of the bastion.
