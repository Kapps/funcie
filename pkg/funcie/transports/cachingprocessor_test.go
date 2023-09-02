package transports_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	. "github.com/Kapps/funcie/pkg/funcie/transports"
	"github.com/Kapps/funcie/pkg/funcie/transports/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCachingMessageProcessor_ForwardRequest(t *testing.T) {
	ctx := context.Background()
	underlying := mocks.NewMessageProcessor(t)
	processor := NewCachingMessageProcessor(underlying)

	msg := &funcie.Message{
		Application: "testApp",
		Kind:        messages.MessageKindForwardRequest,
	}

	underlying.EXPECT().ProcessMessage(ctx, msg).Return(nil, funcie.ErrNoActiveConsumer).Once()

	_, err := processor.ProcessMessage(ctx, msg)
	require.ErrorIs(t, funcie.ErrNoActiveConsumer, err)

	// Test cached result
	_, err = processor.ProcessMessage(ctx, msg)
	require.Equal(t, funcie.ErrNoActiveConsumer, err)
}

func TestCachingMessageProcessor_Register(t *testing.T) {
	ctx := context.Background()
	underlying := mocks.NewMessageProcessor(t)
	processor := NewCachingMessageProcessor(underlying)

	forwardMsg := &funcie.Message{
		Application: "testApp",
		Kind:        messages.MessageKindForwardRequest,
	}

	registerMsg := &funcie.Message{
		Application: "testApp",
		Kind:        messages.MessageKindRegister,
	}

	resp := funcie.NewResponse("resp", nil, nil)

	underlying.EXPECT().ProcessMessage(ctx, forwardMsg).Return(nil, funcie.ErrNoActiveConsumer).Once()

	// First call returns ErrNoActiveConsumer and caches the result
	_, err := processor.ProcessMessage(ctx, forwardMsg)
	require.ErrorIs(t, funcie.ErrNoActiveConsumer, err)

	// Second call returns ErrNoActiveConsumer from cache
	_, err = processor.ProcessMessage(ctx, forwardMsg)
	require.ErrorIs(t, funcie.ErrNoActiveConsumer, err)

	// Register the application
	underlying.EXPECT().ProcessMessage(ctx, registerMsg).Return(resp, nil).Once()
	returnedResp, err := processor.ProcessMessage(ctx, registerMsg)
	require.NoError(t, err)
	require.Equal(t, resp, returnedResp)

	// Now the cached result should be cleared
	underlying.EXPECT().ProcessMessage(ctx, forwardMsg).Return(resp, nil).Once()
	returnedResp, err = processor.ProcessMessage(ctx, forwardMsg)
	require.NoError(t, err)
	require.Equal(t, resp, returnedResp)
}

func TestCachingMessageProcessor_Deregister(t *testing.T) {
	ctx := context.Background()
	underlying := mocks.NewMessageProcessor(t)
	processor := NewCachingMessageProcessor(underlying)

	msg := &funcie.Message{
		Application: "testApp",
		Kind:        messages.MessageKindDeregister,
	}
	resp := funcie.NewResponse("resp", nil, nil)

	underlying.EXPECT().ProcessMessage(ctx, msg).Return(resp, nil).Once()

	returned, err := processor.ProcessMessage(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, resp, returned)
}

func TestCachingMessageProcessor_SuccessfulForwards(t *testing.T) {
	ctx := context.Background()
	underlying := mocks.NewMessageProcessor(t)
	processor := NewCachingMessageProcessor(underlying)

	msg := &funcie.Message{
		Application: "testApp",
		Kind:        messages.MessageKindForwardRequest,
	}
	resp := funcie.NewResponse("resp", nil, nil)

	underlying.EXPECT().ProcessMessage(ctx, msg).Return(resp, nil).Once()

	returned, err := processor.ProcessMessage(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, resp, returned)
}
