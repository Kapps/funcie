package tunnel_test

import (
	"context"
	. "funcie/tunnel"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	redisServer := miniredis.RunT(t)
	consumerClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	defer CloseOrLog(consumerClient)

	producerClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	defer CloseOrLog(producerClient)

	producer := NewRedisPublisher(producerClient, "test-channel")
	consumer := NewRedisConsumer(consumerClient, "test-channel")

	first := NewMessage([]byte("first"), time.Minute)
	second := NewMessage([]byte("second"), time.Minute)

	expectedFirstResponse := NewResponse(first.ID, []byte("resp"))
	expectedSecondResponse := NewResponse(second.ID, []byte("resp"))

	// No consumer yet, so this should fail with ErrNoConsumer
	_, err := producer.Publish(ctx, first)
	require.ErrorIs(t, err, ErrNoActiveConsumer)

	// Now start the consumer, with its own context
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		err := consumer.Consume(consumerCtx, func(ctx context.Context, message *Message) (*Response, error) {
			return NewResponse(message.ID, []byte("resp")), nil
		})
		require.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	// Now we should be able to publish
	resp, err := producer.Publish(ctx, first)
	require.NoError(t, err)
	require.True(t, RoughCompare(expectedFirstResponse, resp))

	// And again
	resp, err = producer.Publish(ctx, second)
	require.NoError(t, err)
	require.True(t, RoughCompare(expectedSecondResponse, resp))

	// Then, cancel the consumer
	cancel()
	time.Sleep(100 * time.Millisecond)

	// And try to publish again
	_, err = producer.Publish(ctx, first)
	require.ErrorIs(t, err, ErrNoActiveConsumer)
}
