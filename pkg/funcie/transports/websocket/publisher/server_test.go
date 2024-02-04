package publisher_test

/*func TestServer_SendMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	message := funcie.NewMessage("app", messages.MessageKindForwardRequest, json.RawMessage(`{"data":"test"}`))
	serializedMessage := json.RawMessage(funcie.MustSerialize(message))
	envelope := &websocket.Envelope{
		Kind: websocket.PayloadKindRequest,
		Data: &serializedMessage,
	}

	connStore := mocks.NewConnectionStore(t)
	exchange := wsMocks.NewExchange(t)
	acceptor := mocks.NewAcceptor(t)
	logger := slog.Default()

	server := publisher.NewServer(connStore, exchange, acceptor, logger)

	t.Run("successful message sending", func(t *testing.T) {
		conn := wsMocks.NewConnection(t)
		connStore.EXPECT().GetConnection("app").
			Return(conn, nil).Once()
		exchange.EXPECT().Send(ctx, conn, message).
			Return(&funcie.Response{}, nil).Once()

		resp, err := server.SendMessage(ctx, message)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("error getting connection", func(t *testing.T) {
		connStore.EXPECT().GetConnection("app").
			Return(nil, errors.New("connection not found")).Once()

		_, err := server.SendMessage(ctx, message)
		require.Error(t, err)
	})

	t.Run("error marshalling message", func(t *testing.T) {
		conn := wsMocks.NewConnection(t)
		connStore.EXPECT().GetConnection("app").
			Return(conn, nil).Once()

		invalidMessage := funcie.NewMessage("app", messages.MessageKindForwardRequest, json.RawMessage(`{invalid json}`))
		_, err := server.SendMessage(ctx, invalidMessage)
		require.Error(t, err)
	})

	t.Run("error writing message", func(t *testing.T) {
		conn := wsMocks.NewConnection(t)
		connStore.EXPECT().GetConnection("app").
			Return(conn, nil).Once()
		conn.EXPECT().Write(ctx, envelope).
			Return(errors.New("write error")).Once()

		_, err := server.SendMessage(ctx, message)
		require.Error(t, err)
	})
}*/

// ReadLoop is covered by the integration tests.
