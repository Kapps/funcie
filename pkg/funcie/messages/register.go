package messages

import (
	"github.com/Kapps/funcie/pkg/funcie"
)

// MessageKindRegistration is a registration request to a server bastion.
const MessageKindRegistration MessageKind = 2

type RegistrationPayload struct {
	// Name is the name of the application.
	Name string `json:"name"`
	// Endpoint is the address to send requests to.
	Endpoint funcie.Endpoint `json:"endpoint"`
}
