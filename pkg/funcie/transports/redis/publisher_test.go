package redis_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	. "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/redis/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisPublisher_Publish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	baseChannelName := faker.Word()
	appId := faker.Word()
	channel := GetChannelNameForApplication(baseChannelName, appId)
	redisClient := mocks.NewPublishClient(t)
	publisher := NewPublisher(redisClient, baseChannelName)

	t.Run("should publish a message to the channel", func(t *testing.T) {
		t.Parallel()

		message := funcie.NewMessage(appId, messages.MessageKindForwardRequest, []byte("hello"))
		serializedMessage, err := json.Marshal(message)
		require.NoError(t, err)

		response := funcie.NewResponse(message.ID, []byte("hello"), nil)
		serializedResponse, err := json.Marshal(response)
		require.NoError(t, err)

		publishResult := redis.NewIntCmd(ctx)
		publishResult.SetVal(1)
		redisClient.On("Publish", ctx, channel, serializedMessage).Return(publishResult)

		responseKey := GetResponseKeyForMessage(baseChannelName, message.ID)
		popResult := redis.NewStringSliceCmd(ctx)
		popResult.SetVal([]string{responseKey, string(serializedResponse)})
		redisClient.On("BRPop", ctx, time.Second, responseKey).Return(popResult)

		resp, err := publisher.Publish(ctx, message)
		require.NoError(t, err)

		require.Equal(t, response, resp)
	})
}
