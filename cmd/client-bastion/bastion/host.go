package bastion

import (
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"io"
	"net/http"
)

// Host is the interface implemented by types that can host a Bastion server.
type Host interface {
	// Listen listens for incoming requests.
	// This function never returns unless an error occurs.
	Listen() error
}

type bastionHost struct {
	httpServer *http.Server
	handler    Handler
}

// NewHost creates a new Host listening on the given address.
func NewHost(address string, handler Handler) Host {
	httpServer := &http.Server{
		Addr: address,
	}
	host := &bastionHost{
		httpServer: httpServer,
		handler:    handler,
	}
	host.setHandlers()

	return host
}

func (h *bastionHost) setHandlers() {
	mux := http.NewServeMux()
	mux.HandleFunc("/dispatch", h.processMessage)
	h.httpServer.Handler = mux
}

func (h *bastionHost) Listen() error {
	if err := h.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (h *bastionHost) processMessage(w http.ResponseWriter, r *http.Request) {
	payloadBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	var message funcie.Message
	err = json.Unmarshal(payloadBytes, &message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid request"))
		return
	}

}

func wrapRequest()
