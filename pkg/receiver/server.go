package receiver

import "net/http"

// ClientBastion allows communication between the client application and the server bastion.
// This allows the logic to be separated from the client application and be language-agnostic.
// In turn, this allows the client application to be more easily written in any language.
type ClientBastion interface {
	// Listen starts listening for requests from the client application.
	Listen()
}

type clientBastion struct {
	server  *http.Server
	handler ClientHandler
}

// NewClientBastion creates a new ClientBastion listening on the given address.
func NewClientBastion(address string, handler ClientHandler) ClientBastion {
	server := &http.Server{
		Addr: address,
	}
	return NewClientBastionWithHTTPServer(server, handler)
}

// NewClientBastionWithHTTPServer creates a new ClientBastion with the given HTTP server.
func NewClientBastionWithHTTPServer(server *http.Server, handler ClientHandler) ClientBastion {
	return &clientBastion{
		server:  server,
		handler: handler,
	}
}

func (c *clientBastion) Listen() {

}
