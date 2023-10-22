package publisher

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie/transports/websocket"
	"log/slog"
	"net/http"
	ws "nhooyr.io/websocket"
)

// Server allows accepting websocket connections.
type Server interface {
	// Listen listens for websocket connections on the given address.
	Listen(ctx context.Context, addr string) error
}

type server struct {
	acceptor Acceptor
	logger   *slog.Logger
}

// NewServer creates a new websocket server with the given acceptor.
func NewServer(acceptor Acceptor, logger *slog.Logger) Server {
	return &server{
		acceptor: acceptor,
		logger:   logger,
	}
}

func (s *server) Listen(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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

			defer func() {
				err := conn.Close(ws.StatusNormalClosure, "")
				slog.Info("Closed connection", "err", err, "remote", r.RemoteAddr)
			}()

			s.logger.Info("Accepted connection", "remote", r.RemoteAddr)
		}),
	}

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listening: %w", err)
	}

	return nil
}

func (s *server) readLoop(ctx context.Context, conn websocket.Connection) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var message interface{}
			err := conn.Read(ctx, &message)
			if err != nil {
				s.logger.Error("Failed to read message", err)
				return
			}

			s.logger.Info("Received message", "message", message)
		}
	}
}
