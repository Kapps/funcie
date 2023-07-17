package provider

import (
	"net/http"
	"net/url"
)

// BastionReceiver represents a receiver that can be used to receive requests from a bastion.
type BastionReceiver interface {
	// Start starts the tunnel. This function never returns.
	// The handler is the handler that will be invoked when a request is received.
	// It is subject to the same restrictions as the handler for the serverless function provider (such as lambda.Start).
	Start()
}

type bastionReceiver struct {
	applicationId   string
	bastionEndpoint url.URL
	server          *http.Server
}

// NewBastionReceiver creates a new BastionReceiver.
func NewBastionReceiver(applicationId string, bastionEndpoint url.URL) BastionReceiver {
	return &bastionReceiver{
		applicationId:   applicationId,
		bastionEndpoint: bastionEndpoint,
		server:          &http.Server{},
	}
}

func (r *bastionReceiver) Start() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(r.handleRequest))

}

func (r *bastionReceiver) handleRequest(w http.ResponseWriter, req *http.Request) {

}
