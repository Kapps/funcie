package redis_test

import (
	"context"
	f "github.com/Kapps/funcie/pkg/funcie"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/redis/mocks"
	"github.com/go-faker/faker/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisConsumer_Consume(t *testing.T) {
	ctx := context.Background()
	redisClient := mocks.NewConsumeClient(t)
	channelName := faker.Word()
	consumer := r.NewConsumerWithClient(redisClient, channelName)

	t.Run("should consume a message from the channel", func(t *testing.T) {
		t.Parallel()

		consumerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		messageChannel := make(chan *redis.Message)
		completedChannel := make(chan struct{})
		pubSub := mocks.NewPubSub(t)

		pubSub.EXPECT().Channel().Return(messageChannel).Times(3)
		pubSub.EXPECT().Close().Return(nil).Once()

		redisClient.EXPECT().Subscribe(consumerCtx, channelName).Return(pubSub).Once()

		go func() {
			defer close(completedChannel)
			err := consumer.Consume(consumerCtx, func(ctx context.Context, message *f.Message) (*f.Response, error) {
				return f.NewResponse(message.ID, []byte("resp"), nil), nil
			})
			require.Equal(t, f.ErrPubSubChannelClosed, err)
		}()

		msg1 := f.NewMessage("app", []byte("msg1"), time.Minute)
		msg2 := f.NewMessage("app", []byte("msg2"), time.Minute)

		redisClient.EXPECT().RPush(
			consumerCtx,
			f.GetResponseKeyForMessage(msg1.ID),
			mock.MatchedBy(RoughCompareMatcherJson(f.NewResponse(msg1.ID, []byte("resp"), nil))),
		).Return(&redis.IntCmd{}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(f.MustSerialize(msg1)),
		})
		require.Empty(t, completedChannel)

		redisClient.EXPECT().RPush(
			consumerCtx,
			f.GetResponseKeyForMessage(msg2.ID),
			mock.MatchedBy(RoughCompareMatcherJson(f.NewResponse(msg2.ID, []byte("resp"), nil))),
		).Return(&redis.IntCmd{}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(f.MustSerialize(msg2)),
		})
		require.Empty(t, completedChannel)

		close(messageChannel)

		ExpectReceiveFromChannel(t, completedChannel)
	})
}
