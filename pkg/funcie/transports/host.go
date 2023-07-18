package transports

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

// Host is the interface implemented by types that can host a Bastion server.
type Host interface {
	// Listen listens for incoming requests.
	// This function never returns unless an error occurs, or Close is called.
	Listen(ctx context.Context) error
	// Close closes the host.
	Close(ctx context.Context) error
}

type bastionHost struct {
	httpServer       *http.Server
	messageProcessor MessageProcessor
}

// NewHost creates a new Host listening on the given address.
func NewHost(address string, messageProcessor MessageProcessor) Host {
	httpServer := &http.Server{
		Addr: address,
	}
	host := &bastionHost{
		httpServer:       httpServer,
		messageProcessor: messageProcessor,
	}
	host.setHandlers()

	return host
}

func (h *bastionHost) setHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dispatch", h.processMessage)
	h.httpServer.Handler = mux
}

func (h *bastionHost) Listen(ctx context.Context) error {
	slog.InfoCtx(ctx, "listening for incoming requests", "address", h.httpServer.Addr)
	if err := h.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (h *bastionHost) Close(ctx context.Context) error {
	slog.Info("closing http server")
	if err := h.httpServer.Close(); err != nil {
		return fmt.Errorf("close http server: %w", err)
	}

	return nil
}

func (h *bastionHost) processMessage(w http.ResponseWriter, r *http.Request) {
	slog.InfoCtx(r.Context(), "received request", "method", r.Method, "url", r.URL)
	payloadBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	slog.DebugCtx(r.Context(), "received payload", "payload", string(payloadBytes))

	var message funcie.Message
	err = json.Unmarshal(payloadBytes, &message)
	if err != nil {
		slog.ErrorCtx(r.Context(), "error unmarshalling message", err, "payload", string(payloadBytes))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("invalid request: %v", err)))
		return
	}

	response, err := h.messageProcessor.ProcessMessage(r.Context(), &message)
	if err != nil {
		slog.ErrorCtx(r.Context(), "error processing message", err, "message", message)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("internal server error: %v", err)))
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		slog.ErrorCtx(r.Context(), "error formatting response", err, "response", response)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("internal server error formatting response: %v", err)))
		return
	}

	_, err = w.Write(responseBytes)
	if err != nil {
		slog.ErrorCtx(r.Context(), "error writing response", err, "response", response)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("internal server error writing response: %v", err)))
		return
	}

	slog.DebugCtx(r.Context(), "sent response", "response", string(responseBytes))
}
