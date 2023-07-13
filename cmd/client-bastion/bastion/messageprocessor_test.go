package bastion_test

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMessageProcessor_ProcessMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := mocks.NewHandler(t)
	processor := bastion.NewMessageProcessor(handler)

	t.Run("registration message", func(t *testing.T) {
		t.Parallel()

		payload := messages.NewRegistrationRequestPayload("app", funcie.MustNewEndpointFromAddress("http://localhost:8080"))
		message := funcie.NewMessageWithPayload("app", messages.MessageKindRegister, *payload)
		response := funcie.NewResponseWithPayload(message.ID, messages.NewRegistrationResponsePayload(uuid.New()), nil)
		handler.EXPECT().Register(ctx, message).Return(response, nil).Once()

		marshaledMessage, err := funcie.MarshalMessagePayload(*message)
		require.NoError(t, err)
		marshaledResponse, err := funcie.MarshalResponsePayload(response)
		require.NoError(t, err)

		resp, err := processor.ProcessMessage(ctx, marshaledMessage)
		require.NoError(t, err)

		RequireEqualResponse(t, response, resp)
	}
}
