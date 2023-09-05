package bastion_test

import (
	"context"
	"encoding/json"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	bastionMocks "github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO: These tests are already unwieldy due to the number of mocks.
// Consider refactoring to use a test suite.

func TestHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	registry := mocks.NewApplicationRegistry(t)
	appClient := bastionMocks.NewApplicationClient(t)
	consumer := mocks.NewConsumer(t)
	hostTranslator := bastionMocks.NewHostTranslator(t)

	hostTranslator.EXPECT().TranslateLocalHostToResolvedHost(ctx, "localhost").Return("localhost", nil)

	handler := bastion.NewHandler(registry, appClient, consumer, hostTranslator)

	endpoint := funcie.MustNewEndpointFromAddress("http://localhost:8080")
	app := funcie.NewApplication("app", endpoint)

	t.Run("should register an application", func(t *testing.T) {
		payload := messages.NewRegistrationRequestPayload(app.Name, app.Endpoint)
		message := funcie.NewMessageWithPayload(app.Name, messages.MessageKindRegister, *payload)

		registry.EXPECT().Register(ctx, app).Return(nil).Once()
		consumer.EXPECT().Subscribe(ctx, app.Name, mock.Anything).Return(nil).Once()

		registered, err := handler.Register(ctx, *message)
		require.NoError(t, err)

		require.NotZero(t, registered.Data.RegistrationId)
	})

	t.Run("should unregister an application", func(t *testing.T) {
		payload := messages.NewDeregistrationRequestPayload(app.Name)
		message := funcie.NewMessageWithPayload(app.Name, messages.MessageKindDeregister, *payload)

		responsePayload := messages.NewDeregistrationResponsePayload()
		expectedResponse := funcie.NewResponseWithPayload(message.ID, responsePayload, nil)

		registry.EXPECT().Unregister(ctx, app.Name).Return(nil).Once()
		consumer.EXPECT().Unsubscribe(ctx, app.Name).Return(nil).Once()

		resp, err := handler.Deregister(ctx, *message)
		require.NoError(t, err)

		RequireEqualResponse(t, expectedResponse, resp)
	})

	t.Run("should wrap an ApplicationNotFound error if the application is not registered", func(t *testing.T) {
		payload := messages.NewDeregistrationRequestPayload(app.Name)
		message := funcie.NewMessageWithPayload(app.Name, messages.MessageKindDeregister, *payload)

		registry.EXPECT().Unregister(ctx, app.Name).Return(funcie.ErrApplicationNotFound).Once()

		resp, err := handler.Deregister(ctx, *message)
		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
		require.Nil(t, resp)
	})

	t.Run("should forward a request to an application", func(t *testing.T) {
		payload := messages.NewForwardRequestPayload(json.RawMessage("{}"))
		request := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, *payload)

		marshaledRequest, err := funcie.MarshalMessagePayload(*request)
		require.NoError(t, err)

		responsePayload := messages.NewForwardRequestResponsePayload(json.RawMessage("{}"))
		response := funcie.NewResponseWithPayload("id", responsePayload, nil)
		marshaledResponse, err := funcie.MarshalResponsePayload(response)
		require.NoError(t, err)

		registry.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()
		appClient.EXPECT().ProcessRequest(ctx, *app, marshaledRequest).Return(marshaledResponse, nil).Once()

		receivedResponse, err := handler.ForwardRequest(ctx, *request)
		require.NoError(t, err)

		RequireEqualResponse(t, response, receivedResponse)
	})

	t.Run("should send a request to an application when consuming a message", func(t *testing.T) {
		registerPayload := messages.NewRegistrationRequestPayload("app", endpoint)
		registerRequest := funcie.NewMessageWithPayload(app.Name, messages.MessageKindRegister, *registerPayload)

		forwardPayload := messages.NewForwardRequestPayload(json.RawMessage("{}"))
		forwardRequest := funcie.NewMessageWithPayload("app", messages.MessageKindForwardRequest, *forwardPayload)

		marshaledForwardRequest, err := funcie.MarshalMessagePayload(*forwardRequest)
		require.NoError(t, err)

		responsePayload := messages.NewForwardRequestResponsePayload(json.RawMessage("{}"))
		response := funcie.NewResponseWithPayload("id", responsePayload, nil)

		marshaledResponse, err := funcie.MarshalResponsePayload(response)
		require.NoError(t, err)

		registry.EXPECT().Register(ctx, app).Return(nil).Once()
		consumer.EXPECT().Subscribe(ctx, "app", mock.Anything).Return(nil).Once()
		registry.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()
		appClient.EXPECT().ProcessRequest(ctx, *app, marshaledForwardRequest).Return(marshaledResponse, nil).Once()

		_, err = handler.Register(ctx, *registerRequest)

		consumeCallback := consumer.Calls[0].Arguments[2].(funcie.Handler)
		resp, err := consumeCallback(ctx, marshaledForwardRequest)
		require.NoError(t, err)

		RequireEqualResponse(t, marshaledResponse, resp)
	})
}
