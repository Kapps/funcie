package receiver_test

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/receiver"
	"github.com/Kapps/funcie/pkg/receiver/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRedisApplicationRegistry(t *testing.T) {
	t.Parallel()

	redisClient := mocks.NewRedisClient(t)
	registry := receiver.NewRedisApplicationRegistry(redisClient)
	ctx := context.Background()

	endpoint := funcie.MustNewEndpointFromAddress("http://localhost:8080")
	app := funcie.NewApplication("app1", endpoint)

	t.Run("should register an application", func(t *testing.T) {
		redisClient.EXPECT().HSet(ctx, "funcie:apps:app1", "endpoint", "http://localhost:8080").
			Return(redis.NewIntCmd(ctx, 1)).Once()

		err := registry.Register(ctx, app)

		require.NoError(t, err)
	})

	t.Run("should unregister an application", func(t *testing.T) {
		redisClient.EXPECT().Del(ctx, "funcie:apps:app1").
			Return(redis.NewIntCmd(ctx, 1)).Once()

		err := registry.Unregister(ctx, "app1")

		require.NoError(t, err)
	})

	t.Run("should get an application", func(t *testing.T) {
		redisClient.EXPECT().HGetAll(ctx, "funcie:apps:app1").
			Return(redis.NewMapStringStringResult(map[string]string{"endpoint": "http://localhost:8080"}, nil)).Once()

		application, err := registry.GetApplication(ctx, "app1")

		require.NoError(t, err)
		require.Equal(t, app, application)
	})

	t.Run("should return an error if the application is not found", func(t *testing.T) {
		redisClient.EXPECT().HGetAll(ctx, "funcie:apps:app1").
			Return(redis.NewMapStringStringResult(map[string]string{}, nil)).Once()

		_, err := registry.GetApplication(ctx, "app1")

		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
	})
}
