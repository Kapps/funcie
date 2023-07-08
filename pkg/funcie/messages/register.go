package messages

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/google/uuid"
)

// MessageKindRegister is a registration request to a server bastion.
const MessageKindRegister funcie.MessageKind = "REGISTER"

// RegistrationMessage is a message containing a registration request.
type RegistrationMessage funcie.MessageBase[RegistrationRequestPayload]

// RegistrationResponse is a message containing a registration response.
type RegistrationResponse funcie.ResponseBase[RegistrationResponsePayload]

type RegistrationRequestPayload struct {
	// Name is the name of the application.
	Name string `json:"name"`
	// Endpoint is the address to send requests to.
	Endpoint funcie.Endpoint `json:"endpoint"`
}

// NewRegistrationRequestPayload creates a new RegistrationRequestPayload with the given name and endpoint.
func NewRegistrationRequestPayload(name string, endpoint funcie.Endpoint) *RegistrationRequestPayload {
	return &RegistrationRequestPayload{
		Name:     name,
		Endpoint: endpoint,
	}
}

// RegistrationResponsePayload is a response to a registration request.
type RegistrationResponsePayload struct {
	// RegistrationId is a unique ID that can be used to deregister the application.
	// For now, this is unused.
	RegistrationId uuid.UUID
}

// NewRegistrationResponsePayload creates a new RegistrationResponsePayload with the given registration ID.
func NewRegistrationResponsePayload(registrationId uuid.UUID) *RegistrationResponsePayload {
	return &RegistrationResponsePayload{
		RegistrationId: registrationId,
	}
}
