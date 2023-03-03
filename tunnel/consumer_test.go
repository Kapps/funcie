package tunnel_test

import (
	"context"
	"encoding/json"
	"errors"
	. "funcie/tunnel"
	"funcie/tunnel/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisConsumer_Consume(t *testing.T) {
	ctx := context.Background()
	redisClient := mocks.NewRedisConsumeClient(t)
	channelName := faker.Word()
	consumer := NewRedisConsumerWithClient(redisClient, channelName)

	t.Run("should consume a message from the channel", func(t *testing.T) {
		t.Parallel()

		consumerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		messageChannel := make(chan *redis.Message)
		completedChannel := make(chan struct{})
		pubSub := mocks.NewRedisPubSub(t)

		pubSub.EXPECT().Channel().Return(messageChannel).Times(3)
		pubSub.EXPECT().Close().Return(nil).Once()

		redisClient.EXPECT().Subscribe(consumerCtx, channelName).Return(pubSub).Once()

		go func() {
			defer close(completedChannel)
			err := consumer.Consume(consumerCtx, func(ctx context.Context, message *Message) (*Response, error) {
				return NewResponse(message.ID, []byte("resp"), nil), nil
			})
			require.Equal(t, ErrPubSubChannelClosed, err)
		}()

		msg1 := NewMessage([]byte("msg1"), time.Minute)
		msg2 := NewMessage([]byte("msg2"), time.Minute)

		redisClient.EXPECT().RPush(
			consumerCtx,
			GetResponseKeyForMessage(msg1.ID),
			mock.MatchedBy(RoughCompareMatcherJson(NewResponse(msg1.ID, []byte("resp"), nil))),
		).Return(&redis.IntCmd{}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(MustSerialize(msg1)),
		})
		require.Empty(t, completedChannel)

		redisClient.EXPECT().RPush(
			consumerCtx,
			GetResponseKeyForMessage(msg2.ID),
			mock.MatchedBy(RoughCompareMatcherJson(NewResponse(msg2.ID, []byte("resp"), nil))),
		).Return(&redis.IntCmd{}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(MustSerialize(msg2)),
		})
		require.Empty(t, completedChannel)

		close(messageChannel)

		ExpectReceiveFromChannel(t, completedChannel)
	})
}

func TestResponseUnmarshal(t *testing.T) {
	t.Parallel()

	response := NewResponse("id", []byte("data"), nil)
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var resp Response
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, response, &resp)
}

func TestResponseUnmarshal_WithError(t *testing.T) {
	t.Parallel()

	response := NewResponse("id", nil, errors.New("error"))
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var resp Response
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, response, &resp)
}
