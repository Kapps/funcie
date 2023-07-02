package bastion_test

import (
	"context"
	"encoding/json"
	. "github.com/Kapps/funcie/cmd/server-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRequestHandler_Dispatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	publisher := mocks.NewPublisher(t)
	ttl := time.Minute * 12
	handler := NewRequestHandler(publisher, ttl)

	t.Run("should publish the request payload to the publisher", func(t *testing.T) {
		payload := json.RawMessage(`{"foo": "bar"}`)
		request := &Request{
			RequestId:         uuid.New().String(),
			Application:       "application",
			Payload:           &payload,
			MessageKind:       funcie.MessageKindDispatch,
			RequestParameters: nil,
		}

		messagePayload, err := request.Payload.MarshalJSON()
		require.NoError(t, err)

		message := funcie.NewMessage("application", request.MessageKind, messagePayload, ttl)

		responsePayload := []byte("response")
		response := funcie.NewResponse(message.ID, responsePayload, nil)

		publisher.EXPECT().Publish(ctx, MessageComparer(message)).Return(response, nil).Once()

		resp, err := handler.Dispatch(ctx, request)
		require.NoError(t, err)

		require.Equal(t, response, resp)
	})
}
