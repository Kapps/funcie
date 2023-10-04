package funcie_test

import (
	"log/slog"
	"os"
	"testing"
)

/*func TestIntegration_Redis(t *testing.T) {
	ctx := context.Background()

	redisServer := miniredis.RunT(t)
	consumerClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	defer CloseOrLog("consumer client", consumerClient)

	producerClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	defer CloseOrLog("producer client", producerClient)

	producer := redistransport.NewPublisher(producerClient, "test-channel")
	consumer := redistransport.NewConsumer(consumerClient, "test-channel")

	first := NewMessage("app", MessageKindRegistration, []byte("first"), time.Minute)
	second := NewMessage("app", MessageKindDispatch, []byte("second"), time.Minute)

	expectedFirstResponse := NewResponse(first.ID, []byte("resp"), nil)
	expectedSecondResponse := NewResponse(second.ID, []byte("resp"), nil)

	// No consumer yet, so this should fail with ErrNoConsumer
	_, err := producer.Publish(ctx, first)
	require.ErrorIs(t, err, ErrNoActiveConsumer)

	// Now start the consumer, with its own context
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := false

	go func() {
		err := consumer.Consume(consumerCtx, func(ctx context.Context, message *Message) (*Response, error) {
			return NewResponse(message.ID, []byte("resp"), nil), nil
		})
		require.True(t, done)
		require.Equal(t, context.Canceled, err)
	}()

	time.Sleep(100 * time.Millisecond)
	log.Println("Starting to publish messages...")

	// Now we should be able to publish
	resp, err := producer.Publish(ctx, first)
	require.NoError(t, err)
	RequireRoughCompare(t, expectedFirstResponse, resp)

	// And again
	resp, err = producer.Publish(ctx, second)
	require.NoError(t, err)
	RequireRoughCompare(t, expectedSecondResponse, resp)

	// Then, cancel the consumer
	log.Println("Cancelling consumer...")
	done = true
	cancel()
	time.Sleep(100 * time.Millisecond)

	// And try to publish again
	_, err = producer.Publish(ctx, first)
	require.ErrorIs(t, err, ErrNoActiveConsumer)
}*/

func TestMain(m *testing.M) {
	programLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     programLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
	programLevel.Set(slog.LevelDebug)
	m.Run()
}
