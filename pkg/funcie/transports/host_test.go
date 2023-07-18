package transports_test

import (
	"bytes"
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestBastionHost_Listen(t *testing.T) {
	ctx := context.Background()
	processor := mocks.NewMessageProcessor(t)
	host := transports.NewHost("localhost:8080", processor)

	go func() {
		err := host.Listen(nil)
		require.ErrorIs(t, err, http.ErrServerClosed)
	}()

	t.Cleanup(func() { _ = host.Close(ctx) })

	client := http.Client{}

	t.Run("success", func(t *testing.T) {

		message := funcie.NewMessage("app", messages.MessageKindRegister, []byte("{}"))
		serialized := funcie.MustSerialize(message)

		response := funcie.NewResponse(message.ID, []byte("{}"), nil)
		processor.EXPECT().ProcessMessage(mock.Anything, message).Return(response, nil).Once()

		resp, err := client.Post("http://localhost:8080/dispatch", "application/json", bytes.NewReader(serialized))
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		responseBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, funcie.MustSerialize(response), responseBytes)
	})

	t.Run("error processing", func(t *testing.T) {
		message := funcie.NewMessage("app", messages.MessageKindRegister, []byte("{}"))
		serialized := funcie.MustSerialize(message)

		processor.EXPECT().ProcessMessage(mock.Anything, message).Return(nil, io.ErrUnexpectedEOF).Once()

		resp, err := client.Post("http://localhost:8080/dispatch", "application/json", bytes.NewReader(serialized))
		require.NoError(t, err)

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		responseBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, "internal server error: unexpected EOF", string(responseBytes))
	})

	t.Run("invalid message", func(t *testing.T) {
		resp, err := client.Post("http://localhost:8080/dispatch", "application/json", bytes.NewReader([]byte("invalid")))
		require.NoError(t, err)

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		responseBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(
			t,
			"invalid request: invalid character 'i' looking for beginning of value",
			string(responseBytes),
		)
	})

	err := host.Close(ctx)
	require.NoError(t, err)
}
