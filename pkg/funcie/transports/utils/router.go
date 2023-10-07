package utils

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"log/slog"
	"sync"
)

var ErrNoHandlerFound = fmt.Errorf("no handler exists for this application")

type ClientHandlerRouter interface {
	AddClientHandler(applicationId string, handler funcie.Handler) error
	RemoveClientHandler(applicationId string) error
	Handle(ctx context.Context, message *funcie.Message) (*funcie.Response, error)
}

func NewClientHandlerRouter() ClientHandlerRouter {
	return &clientHandlerRouter{
		handlers: &sync.Map{},
	}
}

type clientHandlerRouter struct {
	handlers *sync.Map
}

func (h *clientHandlerRouter) AddClientHandler(applicationId string, handler funcie.Handler) error {
	if _, ok := h.handlers.Load(applicationId); ok {
		slog.Warn("overwriting handler for application", "application", applicationId)
	}
	h.handlers.Store(applicationId, handler)
	return nil
}

func (h *clientHandlerRouter) RemoveClientHandler(applicationId string) error {
	if _, deleted := h.handlers.LoadAndDelete(applicationId); !deleted {
		return fmt.Errorf("no handler exists for application %s", applicationId)
	}
	return nil
}

func (h *clientHandlerRouter) Handle(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	handler, ok := h.handlers.Load(message.Application)
	if !ok {
		return nil, fmt.Errorf("application %s not registered: %w", message.Application, ErrNoHandlerFound)
	}
	return handler.(funcie.Handler)(ctx, message)
}

func (h *clientHandlerRouter) ListHandlers() []string {
	var handlers []string
	h.handlers.Range(func(key, value interface{}) bool {
		handlers = append(handlers, key.(string))
		return true
	})
	return handlers
}
