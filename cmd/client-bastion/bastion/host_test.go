package bastion_test

import (
	"bytes"
	"context"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/Kapps/funcie/cmd/client-bastion/bastion/mocks"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestBastionHost_Listen(t *testing.T) {
	ctx := context.Background()
	processor := mocks.NewMessageProcessor(t)
	host := bastion.NewHost("localhost:8080", processor)

	go func() {
		err := host.Listen(nil)
		require.ErrorIs(t, err, http.ErrServerClosed)
	}()

	t.Cleanup(func() { _ = host.Close(ctx) })

	client := http.Client{}

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

	err = host.Close(ctx)
	require.NoError(t, err)
}
