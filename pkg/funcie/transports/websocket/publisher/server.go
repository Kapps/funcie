package publisher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	ws "nhooyr.io/websocket"
	"time"
)

// AcceptTimeout is the timeout for accepting a websocket connection.
// This includes the time it takes to upgrade the connection and to receive the registration message.
const AcceptTimeout = 30 * time.Second

// Server allows accepting websocket connections.
type Server interface {
	// Listen listens for websocket connections on the given address.
	Listen(ctx context.Context, addr string) error
	// Close closes the server.
	Close() error
}

type server struct {
	acceptor Acceptor
	logger   *slog.Logger
	registry Registry
	close    func() error
}

// NewServer creates a new websocket server with the given acceptor.
func NewServer(acceptor Acceptor, registry Registry, logger *slog.Logger) Server {
	return &server{
		acceptor: acceptor,
		logger:   logger,
		registry: registry,
		close:    func() error { return errors.New("not listening") },
	}
}

func (s *server) Listen(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(ctx, AcceptTimeout)
			defer cancel()

			conn, err := s.acceptor.Accept(ctx, rw, r)
			if err != nil {
				// TODO: Is this the right way to close the connection?
				r.Close = true
				_ = r.Body.Close()
				rw.WriteHeader(http.StatusBadRequest)
				_, _ = rw.Write([]byte(fmt.Sprintf("failed to accept connection: %v", err)))

				s.logger.Error("Failed to accept connection", err, "remote", r.RemoteAddr)
				return
			}

			s.logger.Info("Accepted connection", "remote", r.RemoteAddr)

			if err := s.registry.Register(ctx, conn); err != nil {
				s.logger.Error("Failed to register connection", err)
				closeErr := conn.Close(ws.StatusInternalError, "failed to register connection")
				if closeErr != nil {
					s.logger.Error("Failed to close connection", closeErr)
				}

				return
			}

			s.logger.Info("Registered connection", "remote", r.RemoteAddr)
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
	return s.close()
}
