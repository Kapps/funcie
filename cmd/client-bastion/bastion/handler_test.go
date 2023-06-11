package bastion_test

import (
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/mocks"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	registry := mocks.NewApplicationRegistry(t)
	ctx := context.Background()
	app := funcie.NewApplication("name", "endpoint")
	handler := bastion.NewHandler(registry)

	t.Run("should register the handler", func(t *testing.T) {
		registry.EXPECT().Register(ctx, app).Return(nil).Once()

		err := handler.Register(ctx, app)
		require.NoError(t, err)
	})
}

func TestHandler_Unregister(t *testing.T) {
	registry := mocks.NewApplicationRegistry(t)
	ctx := context.Background()
	app := funcie.NewApplication("name", "endpoint")
	handler := bastion.NewHandler(registry)

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
	registry := mocks.NewApplicationRegistry(t)
	ctx := context.Background()

	request := funcie.NewMessage("app", funcie.MessageKindDispatch, []byte("payload"), time.Minute*5)
	response := funcie.NewResponse("id", []byte("response"), nil)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBytes, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		received := funcie.MustDeserialize[*funcie.Message](requestBytes)
		require.Equal(t, request, received)

		respBytes := funcie.MustSerialize(response)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))

	t.Cleanup(srv.Close)

	app := funcie.NewApplication("app", srv.URL)
	handler := bastion.NewHandler(registry)

	registry.EXPECT().GetApplication(ctx, app.Name).Return(app, nil).Once()

	receivedResponse, err := handler.ForwardRequest(ctx, request)
	require.NoError(t, err)

	require.Equal(t, response, receivedResponse)
}
