package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"net/http"
	nhws "nhooyr.io/websocket"
	"time"
)

// AcceptTimeout is the timeout for accepting a websocket connection.
// This includes the time it takes to upgrade the connection and to receive the registration message.
const AcceptTimeout = 30 * time.Second

// Server is responsible for managing all connections and communicating with them.
type Server interface {
	// Listen begins listening for websocket connections on the given address.
	Listen(ctx context.Context, addr string) error
	// Close closes the server.
	Close() error
	// SendMessage sends a message to the application it is intended for.
	// If no connection is found for the given application, ErrNoConnection is returned.
	SendMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

type server struct {
	connStore        ConnectionStore
	responseNotifier ResponseNotifier
	acceptor         Acceptor
	close            func() error
	logger           *slog.Logger
}

// NewServer returns a new websocket server.
// The server will not listen for connections until Listen is called.
func NewServer(
	connStore ConnectionStore,
	responseNotifier ResponseNotifier,
	acceptor Acceptor,
	logger *slog.Logger,
) Server {
	srv := &server{
		connStore:        connStore,
		responseNotifier: responseNotifier,
		acceptor:         acceptor,
		logger:           logger,
		close:            func() error { return errors.New("not listening") },
	}

	return srv
}

func (s *server) Listen(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		//ReadHeaderTimeout: AcceptTimeout,
		//WriteTimeout:      AcceptTimeout,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rootContext := ctx
			ctx, cancel := context.WithTimeout(r.Context(), AcceptTimeout)
			defer cancel()

			s.logger.Info("Received connection", "remote", r.RemoteAddr)

			conn, err := s.acceptor.Accept(ctx, rw, r)
			if err != nil {
				// TODO: Is this the right way to close the connection?
				r.Close = true
				_ = r.Body.Close()
				rw.WriteHeader(http.StatusBadRequest)
				_, _ = rw.Write([]byte(fmt.Sprintf("failed to accept connection: %v", err)))

				s.logger.Error("Failed to accept connection", "error", err, "remote", r.RemoteAddr)
				return
			}

			s.logger.Info("Accepted connection", "remote", r.RemoteAddr)

			go s.readLoop(rootContext, conn)
		}),
	}

	s.close = func() error {
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutting down: %w", err)
		}
		s.close = func() error { return errors.New("not listening") }
		return nil
	}

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listening: %w", err)
	}

	return nil
}

func (s *server) Close() error {
	if err := s.close(); err != nil {
		return fmt.Errorf("closing: %w", err)
	}
	return nil
}

func (s *server) SendMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	s.logger.DebugContext(ctx, "Sending message", "message", message.ID, "application", message.Application)

	conn, err := s.connStore.GetConnection(message.Application)
	if err != nil {
		return nil, fmt.Errorf("getting connection for app %v: %w", message.Application, err)
	}

	s.logger.DebugContext(ctx, "Got connection", "message", message.ID, "application", message.Application)

	payload, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshalling message: %w", err)
	}

	jsonPayload := json.RawMessage(payload)

	err = conn.Write(ctx, &websocket.Envelope{
		Kind: websocket.PayloadKindRequest,
		Data: &jsonPayload,
	})
	if err != nil {
		return nil, fmt.Errorf("writing message: %w", err)
	}

	// TODO: If the connection closes, we should stop waiting early.

	resp, err := s.responseNotifier.WaitForResponse(ctx, message.ID)
	if err != nil {
		return nil, fmt.Errorf("waiting for response: %w", err)
	}

	return resp, nil
}

func (s *server) readLoop(ctx context.Context, conn websocket.Connection) {
	closeConn := func(reason string) {
		if err := conn.Close(nhws.StatusNormalClosure, reason); err != nil {
			s.logger.ErrorContext(ctx, "Error closing connection", "error", err)
		}
	}
	for {
		s.logger.DebugContext(ctx, "Waiting for next message")

		envelope, err := s.readNextMessage(ctx, conn)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error reading message; closing connection", "error", err)
			closeConn("error reading message")
			break
		}

		err = s.processMessage(ctx, envelope, conn)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error processing message; closing connection", "error", err)
			closeConn("error processing message")
			break
		}
	}
}

func (s *server) readNextMessage(ctx context.Context, conn websocket.Connection) (*websocket.Envelope, error) {
	var msg websocket.Envelope
	err := conn.Read(ctx, &msg)
	if err != nil {
		return nil, fmt.Errorf("reading message: %w", err)
	}

	return &msg, nil
}

func (s *server) processMessage(ctx context.Context, envelope *websocket.Envelope, conn websocket.Connection) error {
	switch envelope.Kind {
	case websocket.PayloadKindRequest:
		return s.processRequest(ctx, envelope, conn)
	case websocket.PayloadKindResponse:
		return s.processResponse(ctx, envelope)
	default:
		return fmt.Errorf("invalid message type: %v", envelope.Kind)
	}

	return nil
}

func (s *server) processRequest(ctx context.Context, envelope *websocket.Envelope, conn websocket.Connection) error {
	var msg *funcie.Message
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	s.logger.DebugContext(ctx, "Received request", "message", msg)

	resp, err := s.requestHandler(ctx, conn, msg)
	if err != nil {
		// TODO: Should we send an error response?
		return fmt.Errorf("handling request: %w", err)
	}

	if resp == nil {
		return nil
	}

	responsePayload, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshalling response: %w", err)
	}

	jsonPayload := json.RawMessage(responsePayload)

	err = conn.Write(ctx, &websocket.Envelope{
		Kind: websocket.PayloadKindResponse,
		Data: &jsonPayload,
	})
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

func (s *server) processResponse(ctx context.Context, envelope *websocket.Envelope) error {
	var msg *funcie.Response
	err := json.Unmarshal(*envelope.Data, &msg)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	s.logger.DebugContext(ctx, "Received response", "message", msg)

	s.responseNotifier.Notify(ctx, msg)

	return nil
}

func (s *server) requestHandler(ctx context.Context, conn websocket.Connection, msg *funcie.Message) (*funcie.Response, error) {
	switch msg.Kind {
	case messages.MessageKindRegister:
		s.connStore.RegisterConnection(msg.Application, conn)
		s.logger.InfoContext(ctx, "Registered connection", "application", msg.Application)

		registerPayload := messages.NewRegistrationResponsePayload(uuid.New())
		resp := funcie.NewResponseWithPayload(msg.ID, registerPayload, nil)
		untyped, err := funcie.MarshalResponsePayload(resp)
		if err != nil {
			return nil, fmt.Errorf("marshalling registratoin response payload: %w", err)
		}

		return untyped, nil
	case messages.MessageKindDeregister:
		_, err := s.connStore.UnregisterConnection(msg.Application)
		if err != nil {
			return nil, fmt.Errorf("unregistering connection: %w", err)
		}

		deregisterPayload := messages.NewDeregistrationResponsePayload()
		resp := funcie.NewResponseWithPayload(msg.ID, deregisterPayload, nil)
		untyped, err := funcie.MarshalResponsePayload(resp)
		if err != nil {
			return nil, fmt.Errorf("marshalling deregistration response payload: %w", err)
		}

		s.logger.InfoContext(ctx, "Unregistered connection", "application", msg.Application)
		return untyped, nil
	default:
		return nil, fmt.Errorf("invalid client to server message type: %v", msg.Kind)
	}
}
