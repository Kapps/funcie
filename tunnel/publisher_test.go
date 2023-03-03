package tunnel_test

import (
	"context"
	"encoding/json"
	. "funcie/tunnel"
	"funcie/tunnel/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	t.Parallel()

	t.Run("should return a new message with a unique ID", func(t *testing.T) {
		t.Parallel()

		m1 := NewMessage([]byte("hello"), time.Second)
		m2 := NewMessage([]byte("hello"), time.Second)

		require.NotEqual(t, m1.ID, m2.ID)
	})

	t.Run("should return a new message with the given data", func(t *testing.T) {
		t.Parallel()

		m := NewMessage([]byte("hello"), time.Second)

		require.Equal(t, m.Data, []byte("hello"))
	})
}

func TestRedisPublisher_Publish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	channel := "test-channel"
	redisClient := mocks.NewRedisPublishClient(t)
	publisher := NewRedisPublisher(redisClient, channel)

	t.Run("should publish a message to the channel", func(t *testing.T) {
		t.Parallel()

		message := NewMessage([]byte("hello"), time.Second)
		serializedMessage, err := json.Marshal(message)
		require.NoError(t, err)

		response := NewResponse(message.ID, []byte("hello"), nil)
		serializedResponse, err := json.Marshal(response)
		require.NoError(t, err)

		publishResult := redis.NewIntCmd(ctx)
		publishResult.SetVal(1)
		redisClient.On("Publish", ctx, channel, serializedMessage).Return(publishResult)

		popResult := redis.NewStringSliceCmd(ctx)
		popResult.SetVal([]string{GetResponseKeyForMessage(message.ID), string(serializedResponse)})
		redisClient.On("BRPop", ctx, time.Second, GetResponseKeyForMessage(message.ID)).Return(popResult)

		resp, err := publisher.Publish(ctx, message)
		require.NoError(t, err)

		require.Equal(t, response, resp)
	})
}
