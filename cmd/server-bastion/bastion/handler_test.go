package bastion_test

import (
	"context"
	. "github.com/Kapps/funcie/cmd/server-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRequestHandler(t *testing.T) {
	t.Parallel()

	publisher := mocks.NewPublisher(t)
	handler := NewRequestHandler(publisher)
	require.NotNil(t, handler)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	publisher := mocks.NewPublisher(t)
	handler := NewRequestHandler(publisher)

	resp, err := handler.Register(context.TODO(), messages.RegistrationMessage{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "register unsupported", err.Error())
}

func TestDeregister(t *testing.T) {
	t.Parallel()

	publisher := mocks.NewPublisher(t)
	handler := NewRequestHandler(publisher)

	resp, err := handler.Deregister(context.TODO(), messages.DeregistrationMessage{})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, "deregister unsupported", err.Error())
}

func TestForwardRequest(t *testing.T) {
	t.Parallel()

	publisher := mocks.NewPublisher(t)
	ctx := context.Background()

	handler := NewRequestHandler(publisher)

	forwardPayload := messages.NewForwardRequestPayload([]byte("{}"))
	forwardMessage := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, *forwardPayload)

	marshaledForwardMessage, err := funcie.MarshalMessagePayload(*forwardMessage)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		responsePayload := messages.NewForwardRequestResponsePayload([]byte("{}"))
		response := funcie.NewResponseWithPayload(forwardMessage.ID, responsePayload, nil)

		marshaledResponse, err := funcie.MarshalResponsePayload(response)
		require.NoError(t, err)

		publisher.EXPECT().Publish(ctx, marshaledForwardMessage).Return(marshaledResponse, nil).Once()

		resp, err := handler.ForwardRequest(ctx, *forwardMessage)
		require.NotNil(t, resp)
		require.NoError(t, err)
		require.Equal(t, *response, *resp)
	})

	t.Run("no active consumer", func(t *testing.T) {
		response := funcie.NewResponseWithPayload[messages.ForwardRequestResponsePayload](
			forwardMessage.ID, nil, funcie.ErrNoActiveConsumer,
		)

		publisher.EXPECT().Publish(ctx, marshaledForwardMessage).Return(nil, funcie.ErrNoActiveConsumer).Once()

		resp, err := handler.ForwardRequest(ctx, *forwardMessage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, response, resp)
	})

	t.Run("application not found", func(t *testing.T) {
		// We expect the same response as for no active consumer
		response := funcie.NewResponseWithPayload[messages.ForwardRequestResponsePayload](
			forwardMessage.ID, nil, funcie.ErrNoActiveConsumer,
		)

		publisher.EXPECT().Publish(ctx, marshaledForwardMessage).Return(nil, funcie.ErrApplicationNotFound).Once()

		resp, err := handler.ForwardRequest(ctx, *forwardMessage)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, response, resp)
	})
}
