package redis_test

/*
import (
	"context"
	"encoding/json"
	f "github.com/Kapps/funcie/pkg/funcie"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/redis/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisPublisher_Publish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	channel := "test-channel"
	redisClient := mocks.NewPublishClient(t)
	publisher := r.NewPublisher(redisClient, channel)

	t.Run("should publish a message to the channel", func(t *testing.T) {
		t.Parallel()

		message := f.NewMessage("app", f.MessageKindDispatch, []byte("hello"), time.Second)
		serializedMessage, err := json.Marshal(message)
		require.NoError(t, err)

		response := f.NewResponse(message.ID, []byte("hello"), nil)
		serializedResponse, err := json.Marshal(response)
		require.NoError(t, err)

		publishResult := redis.NewIntCmd(ctx)
		publishResult.SetVal(1)
		redisClient.On("Publish", ctx, channel, serializedMessage).Return(publishResult)

		popResult := redis.NewStringSliceCmd(ctx)
		popResult.SetVal([]string{f.GetResponseKeyForMessage(message.ID), string(serializedResponse)})
		redisClient.On("BRPop", ctx, time.Second, f.GetResponseKeyForMessage(message.ID)).Return(popResult)

		resp, err := publisher.Publish(ctx, message)
		require.NoError(t, err)

		require.Equal(t, response, resp)
	})
}
*/
