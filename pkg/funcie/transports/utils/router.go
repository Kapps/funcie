package utils

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"log/slog"
)

var ErrNoHandlerFound = fmt.Errorf("no handler exists for this application")

type ClientHandlerRouter interface {
	AddClientHandler(applicationId string, handler funcie.Handler) error
	RemoveClientHandler(applicationId string) error
	Handle(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

func NewClientHandlerRouter() ClientHandlerRouter {
	return &clientHandlerRouter{
		handlers: map[string]funcie.Handler{},
	}
}

type clientHandlerRouter struct {
	handlers map[string]funcie.Handler
}

func (h *clientHandlerRouter) AddClientHandler(applicationId string, handler funcie.Handler) error {
	if _, ok := h.handlers[applicationId]; ok {
		slog.Warn("overwriting handler for application", "application", applicationId)
	}
	h.handlers[applicationId] = handler
	return nil
}

func (h *clientHandlerRouter) RemoveClientHandler(applicationId string) error {
	if _, ok := h.handlers[applicationId]; !ok {
		return fmt.Errorf("no handler exists for application %s", applicationId)
	}
	delete(h.handlers, applicationId)
	return nil
}

func (h *clientHandlerRouter) Handle(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	handler, ok := h.handlers[message.Application]
	if !ok {
		return nil, fmt.Errorf("application %s not registered: %w", message.Application, ErrNoHandlerFound)
	}
	return handler(ctx, message)
}

func (h *clientHandlerRouter) ListHandlers() []string {
	var handlers = make([]string, 0, len(h.handlers))
	for handler := range h.handlers {
		handlers = append(handlers, handler)
	}
	return handlers
}
