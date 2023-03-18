package receiver

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
)

type ClientHandler interface {
	RegisterApplication(ctx context.Context, application *funcie.Application) error
	UnregisterApplication(ctx context.Context, applicationName string) error
	DispatchRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error)
}

type clientHandler struct {
	registry funcie.ApplicationRegistry
}

func NewClientHandler(registry funcie.ApplicationRegistry) ClientHandler {
	return &clientHandler{
		registry: registry,
	}
}

func (h *clientHandler) RegisterApplication(ctx context.Context, application *funcie.Application) error {
	return h.registry.Register(ctx, application)
}


func (h *clientHandler) UnregisterApplication(ctx context.Context, applicationName string) error {
	return h.registry.Unregister(ctx, applicationName)
}

func (h *clientHandler) DispatchRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error) {
	application, err := h.registry.GetApplication(ctx, request.)
	if err != nil {
		return nil, fmt.Errorf("failed to get application %s: %w", request.ApplicationName, err)
	}

	return application.Endpoint.Dispatch(ctx, request)
}