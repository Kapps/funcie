package bastion_test

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	funcieMocks "github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestForwardingApplicationRegistry_Integration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := funcie.NewApplication("name", funcie.MustNewEndpointFromAddress("http://localhost:8080"))

	underlying := funcieMocks.NewApplicationRegistry(t)
	publisher := funcieMocks.NewPublisher(t)
	registry := bastion.NewForwardingApplicationRegistry(underlying, publisher)

	t.Run("should forward registration requests", func(t *testing.T) {
		t.Parallel()

		message := funcie.NewMessageWithPayload(
			app.Name,
			messages.MessageKindRegister,
			messages.NewRegistrationRequestPayload(app.Name, app.Endpoint),
			time.Minute*2,
		)

		response := funcie.NewResponseWithPayload(
			message.ID,
			messages.NewRegistrationResponsePayload(uuid.New()),
			nil,
		)
		underlying.EXPECT().Register(ctx, app).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, MessageComparer(message)).Return(response, nil).Once()

		err := registry.Register(ctx, app)
		require.NoError(t, err)
	})

	t.Run("should forward deregistration requests", func(t *testing.T) {
		t.Parallel()

		message := funcie.NewMessageWithPayload(
			app.Name,
			messages.MessageKindDeregister,
			messages.NewDeregistrationRequestPayload(app.Name),
			time.Minute*2,
		)
		response := funcie.NewResponseWithPayload(
			message.ID,
			messages.NewDeregistrationRequestPayload(app.Name),
			nil,
		)
		underlying.EXPECT().Unregister(ctx, app.Name).Return(nil).Once()
		publisher.EXPECT().Publish(ctx, MessageComparer(message)).Return(response, nil).Once()

		err := registry.Unregister(ctx, app.Name)
		require.NoError(t, err)
	})

	t.Run("should pass through application retrieval", func(t *testing.T) {
		t.Parallel()

		underlying.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()

		application, err := registry.GetApplication(ctx, app.Name)
		require.NoError(t, err)
		require.Equal(t, app, application)
	})
}