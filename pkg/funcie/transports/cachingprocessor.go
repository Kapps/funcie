package transports

import (
	"context"
	"errors"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"golang.org/x/exp/slog"
	"sync"
	"time"
)

type cachedEntry struct {
	timestamp time.Time
}

type cachingMessageProcessor struct {
	underlyingProcessor MessageProcessor
	noConsumerCache     sync.Map
}

// NOTE: Currently this is not being used.
// The problem is that client bastions have no way to communicate with a server bastion.
// So the unavailable cache has no way to be cleared beyond waiting the expiration time.
// Eventually we'd probably want some sort of control plane for things like authentication, which
// could then be used to clear the cache upon an application registration.

// Meanwhile, it's not particularly expensive to go from Server -> Client -> Server.

// NewCachingMessageProcessor creates a new caching MessageProcessor, forwarding requests to the underlying processor.
// If the underlying processor returns ErrNoConsumerFound, that result is cached for a minute or until registered.
func NewCachingMessageProcessor(underlyingProcessor MessageProcessor) MessageProcessor {
	return &cachingMessageProcessor{
		underlyingProcessor: underlyingProcessor,
	}
}

func (cp *cachingMessageProcessor) ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	switch message.Kind {
	case messages.MessageKindForwardRequest:
		return cp.handleForwardRequest(ctx, message)
	case messages.MessageKindRegister:
		return cp.handleRegister(ctx, message)
	case messages.MessageKindDeregister:
		return cp.underlyingProcessor.ProcessMessage(ctx, message) // Deregister can directly use the underlying processor
	default:
		return nil, ErrUnknownMessageKind
	}
}

func (cp *cachingMessageProcessor) handleForwardRequest(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	value, found := cp.noConsumerCache.Load(message.Application)
	if found {
		entry := value.(cachedEntry)
		if time.Since(entry.timestamp) < time.Minute {
			slog.DebugCtx(ctx, "no consumer found, cached", "application", message.Application)
			return nil, funcie.ErrNoActiveConsumer
		} else {
			slog.DebugCtx(ctx, "no consumer found, cache expired", "application", message.Application)
			cp.noConsumerCache.Delete(message.Application)
		}
	}

	resp, err := cp.underlyingProcessor.ProcessMessage(ctx, message)
	if errors.Is(err, funcie.ErrNoActiveConsumer) {
		// No consumer because client bastion is unreachable.
		slog.DebugCtx(ctx, "no consumer found (client bastion unresponsive?), caching for a minute", "application", message.Application)
		cp.noConsumerCache.Store(message.Application, cachedEntry{timestamp: time.Now()})
	} else if err == nil && resp.Error != nil && resp.Error.Error() == funcie.ErrNoActiveConsumer.Error() {
		// Client bastion was reachable, and it responded with no consumer.
		slog.DebugCtx(ctx, "no consumer found (negative response), caching for a minute", "application", message.Application)
		cp.noConsumerCache.Store(message.Application, cachedEntry{timestamp: time.Now()})
	}

	return resp, err
}

func (cp *cachingMessageProcessor) handleRegister(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	cp.noConsumerCache.Delete(message.Application)
	return cp.underlyingProcessor.ProcessMessage(ctx, message)
}
