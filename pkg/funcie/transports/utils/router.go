package utils

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
)

type ClientHandlerRouter interface {
	AddClientHandler(handlerType string, handler funcie.Handler) error
	RemoveClientHandler(handlerType string) error
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

func (h *clientHandlerRouter) AddClientHandler(handlerType string, handler funcie.Handler) error {
	if _, ok := h.handlers[handlerType]; ok {
		return fmt.Errorf("handler already exists for type %s", handlerType)
	}
	h.handlers[handlerType] = handler
	return nil
}

func (h *clientHandlerRouter) RemoveClientHandler(handlerType string) error {
	if _, ok := h.handlers[handlerType]; !ok {
		return fmt.Errorf("no handler exists for type %s", handlerType)
	}
	delete(h.handlers, handlerType)
	return nil
}

func (h *clientHandlerRouter) Handle(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	handler, ok := h.handlers[message.Application]
	if !ok {
		return nil, fmt.Errorf("no handler exists for type %s", message.Application)
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
