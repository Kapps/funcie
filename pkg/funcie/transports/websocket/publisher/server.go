package publisher

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"log/slog"
	"net"
	"net/http"
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
}

type server struct {
	connStore ConnectionStore
	exchange  websocket.Exchange
	acceptor  Acceptor
	close     func() error
	logger    *slog.Logger
}

// NewServer returns a new websocket server.
// The server will not listen for connections until Listen is called.
func NewServer(
	connStore ConnectionStore,
	exchange websocket.Exchange,
	acceptor Acceptor,
	logger *slog.Logger,
) Server {
	srv := &server{
		connStore: connStore,
		acceptor:  acceptor,
		logger:    logger,
		exchange:  exchange,
		close:     func() error { return errors.New("not listening") },
	}

	return srv
}

func (s *server) Listen(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: AcceptTimeout,
		WriteTimeout:      AcceptTimeout,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rootContext := ctx
			ctx, cancel := context.WithTimeout(r.Context(), AcceptTimeout)
			defer cancel()

			closeRequest := func() {
				// TODO: Is this the right way to close the connection?
				r.Close = true
				_ = r.Body.Close()
				rw.WriteHeader(http.StatusBadRequest)
				_, _ = rw.Write([]byte("failed to accept connection"))
			}

			s.logger.Info("Received connection", "remote", r.RemoteAddr)

			conn, err := s.acceptor.Accept(ctx, rw, r)
			if err != nil {
				closeRequest()
				s.logger.Error("Failed to accept connection", "error", err, "remote", r.RemoteAddr)
				return
			}

			s.logger.Info("Accepted connection", "remote", r.RemoteAddr)

			if err := s.exchange.RegisterConnection(rootContext, conn); err != nil {
				closeRequest()
				s.logger.Error("Failed to register connection", "error", err, "remote", r.RemoteAddr)
				return
			}
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
