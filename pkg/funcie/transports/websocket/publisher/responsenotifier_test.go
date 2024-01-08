package publisher_test

import (
	"context"
	"errors"
	"github.com/Kapps/funcie/pkg/funcie"
	. "github.com/Kapps/funcie/pkg/funcie/transports/websocket/publisher"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestResponseNotifier_NotifyBeforeWait(t *testing.T) {
	t.Parallel()

	respNotifier := NewResponseNotifier()
	response := &funcie.Response{ID: "1"}
	ctx := context.Background()

	go respNotifier.Notify(ctx, response)

	resp, err := respNotifier.WaitForResponse(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, response, resp)
}

func TestResponseNotifier_WaitBeforeNotify(t *testing.T) {
	t.Parallel()

	respNotifier := NewResponseNotifier()
	response := &funcie.Response{ID: "1"}
	ctx := context.Background()

	var receivedResponse *funcie.Response
	var err error
	done := make(chan struct{})

	go func() {
		receivedResponse, err = respNotifier.WaitForResponse(ctx, "1")
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)

	respNotifier.Notify(ctx, response)
	<-done

	require.NoError(t, err)
	require.Equal(t, response, receivedResponse)
}

func TestResponseNotifier_ContextCancellation(t *testing.T) {
	t.Parallel()

	respNotifier := NewResponseNotifier()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := respNotifier.WaitForResponse(ctx, "1")
	require.Nil(t, resp)
	require.True(t, errors.Is(err, context.Canceled))
}

func TestResponseNotifier_NotifyPanic(t *testing.T) {
	respNotifier := NewResponseNotifier()
	response := &funcie.Response{ID: "1"}
	ctx := context.Background()

	respNotifier.Notify(ctx, response)

	require.Panics(t, func() {
		respNotifier.Notify(ctx, response)
	})
}

func TestResponseNotifier_WaitForNonExistingMessageID(t *testing.T) {
	respNotifier := NewResponseNotifier()
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancelFn()

	resp, err := respNotifier.WaitForResponse(ctx, "non-existing")
	require.Nil(t, resp)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestResponseNotifier_ConcurrentCalls(t *testing.T) {
	respNotifier := NewResponseNotifier()
	ctx := context.Background()

	go respNotifier.Notify(ctx, &funcie.Response{ID: "1"})
	go respNotifier.Notify(ctx, &funcie.Response{ID: "2"})

	resp1, err1 := respNotifier.WaitForResponse(ctx, "1")
	resp2, err2 := respNotifier.WaitForResponse(ctx, "2")

	require.NoError(t, err1)
	require.Equal(t, "1", resp1.ID)
	require.NoError(t, err2)
	require.Equal(t, "2", resp2.ID)
}
