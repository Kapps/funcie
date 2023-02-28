package tunnel_test

import (
	"context"
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

		messageChannel := make(chan *redis.Message)
		completedChannel := make(chan struct{})
		pubSub := mocks.NewRedisPubSub(t)
		pubSub.EXPECT().Channel().Return(messageChannel).Times(3)
		pubSub.EXPECT().Close().Return(nil).Once()

		redisClient.EXPECT().Subscribe(ctx, channelName).Return(pubSub).Once()

		go func() {
			defer close(completedChannel)
			err := consumer.Consume(ctx, func(ctx context.Context, message *Message) (*Response, error) {
				return NewResponse(message.ID, []byte("resp")), nil
			})
			require.NoError(t, err)
		}()

		msg1 := NewMessage([]byte("msg1"), time.Minute)
		msg2 := NewMessage([]byte("msg2"), time.Minute)

		redisClient.EXPECT().RPush(
			ctx,
			GetResponseKeyForMessage(msg1.ID),
			mock.MatchedBy(RoughCompareMatcherJson(NewResponse(msg1.ID, []byte("resp")))),
		).Return(&redis.IntCmd{}).Once()

		messageChannel <- &redis.Message{
			Payload: string(MustSerialize(msg1)),
		}
		require.Empty(t, completedChannel)

		redisClient.EXPECT().RPush(
			ctx,
			GetResponseKeyForMessage(msg2.ID),
			mock.MatchedBy(RoughCompareMatcherJson(NewResponse(msg2.ID, []byte("resp")))),
		).Return(&redis.IntCmd{}).Once()

		messageChannel <- &redis.Message{
			Payload: string(MustSerialize(msg2)),
		}
		require.Empty(t, completedChannel)

		close(messageChannel)

		<-completedChannel
	})
}
