package bastion_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	bastionMocks "github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/google/uuid"
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
	payload := messages.NewRegistrationRequestPayload(app.Name, app.Endpoint)
	message := funcie.NewMessageWithPayload(app.Name, messages.MessageKindRegister, *payload, time.Minute*5)

	t.Run("should register the handler", func(t *testing.T) {
		registry.EXPECT().Register(ctx, app).Return(nil).Once()

		registered, err := handler.Register(ctx, *message)
		require.NoError(t, err)

		require.NotZero(t, registered.Data.RegistrationId)
	})
}

func TestHandler_Unregister(t *testing.T) {
	ctx := context.Background()
	registry := mocks.NewApplicationRegistry(t)
	appClient := bastionMocks.NewApplicationClient(t)
	endpoint, err := funcie.NewEndpointFromAddress("http://localhost:8080")
	require.NoError(t, err)

	app := funcie.NewApplication("name", endpoint)
	payload := messages.NewDeregistrationRequestPayload(app.Name)
	message := funcie.NewMessageWithPayload(app.Name, messages.MessageKindDeregister, *payload, time.Minute*5)
	handler := bastion.NewHandler(registry, appClient)
	responsePayload := messages.NewDeregistrationRequestPayload(uuid.New().String())
	expectedResponse := funcie.NewResponseWithPayload(message.ID, responsePayload, nil)

	t.Run("should unregister the handler", func(t *testing.T) {
		t.Parallel()

		registry.EXPECT().Unregister(ctx, app.Name).Return(nil).Once()

		resp, err := handler.Deregister(ctx, *message)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, resp)
	})

	t.Run("should wrap an ApplicationNotFound error if the application is not registered", func(t *testing.T) {
		t.Parallel()

		registry.EXPECT().Unregister(ctx, app.Name).Return(funcie.ErrApplicationNotFound).Once()

		resp, err := handler.Deregister(ctx, *message)
		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
		require.Nil(t, resp)
	})
}

func TestHandler_ForwardRequest(t *testing.T) {
	ctx := context.Background()
	registry := mocks.NewApplicationRegistry(t)
	appClient := bastionMocks.NewApplicationClient(t)
	endpoint, err := funcie.NewEndpointFromAddress("http://localhost:8080")
	payload := messages.NewForwardRequestPayload(json.RawMessage("{}"))
	require.NoError(t, err)

	request := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, *payload, time.Minute*5)
	require.NoError(t, err)

	response := funcie.NewResponse("id", []byte("response"), nil)

	app := funcie.NewApplication("app", endpoint)
	handler := bastion.NewHandler(registry, appClient)

	registry.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()
	appClient.EXPECT().ProcessRequest(ctx, *app, request).Return(response, nil).Once()

	receivedResponse, err := handler.ForwardRequest(ctx, *request)
	require.NoError(t, err)

	require.Equal(t, response, receivedResponse)
}
