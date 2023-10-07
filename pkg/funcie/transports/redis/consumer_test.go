package redis_test

import (
	"context"
	f "github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	r "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/Kapps/funcie/pkg/funcie/transports/redis/mocks"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	utilMocks "github.com/Kapps/funcie/pkg/funcie/transports/utils/mocks"
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
	baseChannelName := faker.Word()
	appId := faker.Word()
	channelName := r.GetChannelNameForApplication(baseChannelName, appId)
	router := utilMocks.NewClientHandlerRouter(t)
	consumer := r.NewConsumerWithClient(redisClient, baseChannelName, router)

	t.Run("should consume a message from the channel", func(t *testing.T) {
		consumerCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		messageChannel := make(chan *redis.Message)
		completedChannel := make(chan struct{})
		pubSub := mocks.NewPubSub(t)

		// First we connect, which is done by subbing to the base channel and waiting for a response.
		redisClient.EXPECT().Subscribe(ctx, baseChannelName).Return(pubSub).Once()
		pubSub.EXPECT().Receive(ctx).Return(&redis.Subscription{
			Channel: "foo",
			Count:   1,
		}, nil).Once()

		require.NoError(t, consumer.Connect(ctx))

		// Then we connect, which will open the channel and start consuming.
		pubSub.EXPECT().Channel().Return(messageChannel).Times(4)
		pubSub.EXPECT().Close().Return(nil).Once()

		go func() {
			defer close(completedChannel)
			err := consumer.Consume(consumerCtx)
			require.Equal(t, f.ErrPubSubChannelClosed, err)
		}()

		time.Sleep(50 * time.Millisecond)

		msg1 := f.NewMessage(appId, messages.MessageKindForwardRequest, []byte("\"msg1\""))
		msg2 := f.NewMessage(appId, messages.MessageKindForwardRequest, []byte("\"msg2\""))

		// If no handler, receiving a message should unsubscribe and close the channel.
		router.EXPECT().Handle(consumerCtx, msg1).
			RunAndReturn(func(ctx context.Context, message *f.Message) (*f.Response, error) {
				completedChannel <- struct{}{}
				return nil, utils.ErrNoHandlerFound
			}).Once()
		router.EXPECT().RemoveClientHandler(appId).Return(nil).Once()
		pubSub.EXPECT().Unsubscribe(consumerCtx, channelName).Return(nil).Once()
		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(f.MustSerialize(msg1)),
		})

		<-completedChannel

		// Now subscribe, and then try again.

		resp1 := f.NewResponse(msg1.ID, []byte("\"resp1\""), nil)
		resp2 := f.NewResponse(msg2.ID, []byte("\"resp2\""), nil)

		router.EXPECT().AddClientHandler(appId, mock.Anything).Return(nil).Once()
		pubSub.EXPECT().Subscribe(ctx, channelName).Return(nil).Once()

		require.NoError(t, consumer.Subscribe(ctx, appId, f.Handler(nil)))

		redisClient.EXPECT().RPush(
			consumerCtx,
			r.GetResponseKeyForMessage(baseChannelName, msg1.ID),
			mock.MatchedBy(RoughCompareMatcherJson(resp1)),
		).Return(&redis.IntCmd{}).Once()

		router.EXPECT().Handle(
			consumerCtx,
			mock.MatchedBy(RoughCompareMatcher(msg1)),
		).RunAndReturn(func(ctx context.Context, message *f.Message) (*f.Response, error) {
			completedChannel <- struct{}{}
			return resp1, nil
		}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(f.MustSerialize(msg1)),
		})
		<-completedChannel

		redisClient.EXPECT().RPush(
			consumerCtx,
			r.GetResponseKeyForMessage(baseChannelName, msg2.ID),
			mock.MatchedBy(RoughCompareMatcherJson(resp2)),
		).Return(&redis.IntCmd{}).Once()

		router.EXPECT().Handle(
			consumerCtx,
			mock.MatchedBy(RoughCompareMatcher(msg2)),
		).RunAndReturn(func(ctx context.Context, message *f.Message) (*f.Response, error) {
			completedChannel <- struct{}{}
			return resp2, nil
		}).Once()

		ExpectSendToChannel(t, messageChannel, &redis.Message{
			Payload: string(f.MustSerialize(msg2)),
		})
		<-completedChannel

		close(messageChannel)

		ExpectReceiveFromChannel(t, completedChannel)
	})
}
