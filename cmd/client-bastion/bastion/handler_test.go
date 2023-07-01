package bastion_test

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	bastionMocks "github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	registry := mocks.NewApplicationRegistry(t)
	ctx := context.Background()
	app := funcie.NewApplication("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))
	appClient := bastionMocks.NewApplicationClient(t)
	handler := bastion.NewHandler(registry, appClient)

	t.Run("should register the handler", func(t *testing.T) {
		registry.EXPECT().Register(ctx, app).Return(nil).Once()

		err := handler.Register(ctx, app)
		require.NoError(t, err)
	})
}

func TestHandler_Unregister(t *testing.T) {
	ctx := context.Background()
	registry := mocks.NewApplicationRegistry(t)
	appClient := bastionMocks.NewApplicationClient(t)
	endpoint, err := funcie.NewEndpointFromAddress("http://localhost:8080")
	require.NoError(t, err)

	app := funcie.NewApplication("name", endpoint)
	handler := bastion.NewHandler(registry, appClient)

	t.Run("should unregister the handler", func(t *testing.T) {
		t.Parallel()

		registry.EXPECT().Unregister(ctx, app.Name).Return(nil).Once()

		err := handler.Unregister(ctx, app.Name)
		require.NoError(t, err)
	})

	t.Run("should wrap an ApplicationNotFound error if the application is not registered", func(t *testing.T) {
		t.Parallel()

		registry.EXPECT().Unregister(ctx, app.Name).Return(funcie.ErrApplicationNotFound).Once()

		err := handler.Unregister(ctx, app.Name)
		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
	})
}

func TestHandler_ForwardRequest(t *testing.T) {
	ctx := context.Background()
	registry := mocks.NewApplicationRegistry(t)
	appClient := bastionMocks.NewApplicationClient(t)
	endpoint, err := funcie.NewEndpointFromAddress("http://localhost:8080")
	require.NoError(t, err)

	request := funcie.NewMessage("app", funcie.MessageKindDispatch, []byte("payload"), time.Minute*5)
	response := funcie.NewResponse("id", []byte("response"), nil)

	app := funcie.NewApplication("app", endpoint)
	handler := bastion.NewHandler(registry, appClient)

	registry.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()
	appClient.EXPECT().ProcessRequest(ctx, *app, request).Return(response, nil).Once()

	receivedResponse, err := handler.ForwardRequest(ctx, request)
	require.NoError(t, err)

	require.Equal(t, response, receivedResponse)
}
