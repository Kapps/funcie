package ws_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/consumer"
	"github.com/Kapps/funcie/pkg/funcie/transports/ws/publisher"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWS_End2End_Subscribe(t *testing.T) {
	t.Parallel()

	go func() {
		err := publisher.Listen(8086)
		require.NoError(t, err)
	}()

	ctx, cancel := context.WithCancel(context.Background())

	client := consumer.NewConsumer("ws://localhost:8086")
	err := client.Connect(ctx)
	require.NoError(t, err)

	err = client.Subscribe(ctx, "channelName", func(ctx context.Context, msg *funcie.Message) (*funcie.Response, error) {
		return nil, nil
	})
	require.NoError(t, err)

	err = client.Unsubscribe(ctx, "channelName")
	require.NoError(t, err)

	time.Sleep(500 * time.Millisecond)
	cancel()
}
