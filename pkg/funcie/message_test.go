package funcie_test

import (
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUnmarshalPayload(t *testing.T) {
	regPayload := messages.NewRegistrationRequestPayload("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))
	serializedPayload := funcie.MustSerialize(regPayload)

	marshaledMessage := funcie.NewMessage("name", messages.MessageKindRegister, serializedPayload, time.Minute)

	unmarshaledMessage, err := funcie.UnmarshalMessagePayload[messages.RegistrationMessage](marshaledMessage)
	require.NoError(t, err)

	require.Equal(t, *regPayload, unmarshaledMessage.Payload)
}

func TestMarshalPayload(t *testing.T) {
	regPayload := messages.NewRegistrationRequestPayload("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))

	unmarshaledMessage := funcie.NewMessageWithPayload("name", messages.MessageKindRegister, regPayload, time.Minute)

	marshaledMessage, err := funcie.MarshalMessagePayload(*unmarshaledMessage)
	require.NoError(t, err)

	require.Equal(t, json.RawMessage(funcie.MustSerialize(regPayload)), marshaledMessage.Payload)
}
