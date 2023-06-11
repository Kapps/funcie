package bastion

import "net/http"

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
