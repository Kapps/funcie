package messages

import (
	"github.com/Kapps/funcie/pkg/funcie"
)

// MessageKindDeregister is a deregistration request to a server bastion.
const MessageKindDeregister funcie.MessageKind = "DEREGISTER"

// DeregistrationMessage is a message containing a deregistration request.
type DeregistrationMessage = funcie.MessageBase[DeregistrationRequestPayload]

// DeregistrationResponse is a message containing a deregistration response.
type DeregistrationResponse = funcie.ResponseBase[DeregistrationResponsePayload]

// DeregistrationRequestPayload is a deregistration request.
type DeregistrationRequestPayload struct {
	// Name is the name of the application.
	Name string `json:"name"`
}

// NewDeregistrationRequestPayload creates a new DeregistrationRequestPayload with the given name.
func NewDeregistrationRequestPayload(name string) *DeregistrationRequestPayload {
	return &DeregistrationRequestPayload{
		Name: name,
	}
}

// DeregistrationResponsePayload is a response to a deregistration request.
type DeregistrationResponsePayload struct {
}
