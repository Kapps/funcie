package publisher

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"sync"
)

// ResponseNotifier notifies a publisher that a response has been received for a message.
type ResponseNotifier interface {
	// Notify notifies the notifier that a response has been received.
	// This method is non-blocking and may only be called once for each message ID.
	Notify(ctx context.Context, resp *funcie.Response)

	// WaitForResponse waits for a response for the given message ID.
	// This method is blocking and may only be called once for each message ID.
	WaitForResponse(ctx context.Context, messageId string) (*funcie.Response, error)
}

type responseNotifier struct {
	responses map[string]chan *funcie.Response
	mapLock   sync.Mutex
}

// NewResponseNotifier creates a new ResponseNotifier.
func NewResponseNotifier() ResponseNotifier {
	return &responseNotifier{
		responses: make(map[string]chan *funcie.Response),
	}
}

func (r *responseNotifier) Notify(ctx context.Context, resp *funcie.Response) {
	ch := r.getResponseChannel(ctx, resp.ID)

	// This won't block because we only call this method once for each message ID.
	// But let's be safe anyway and panic if it does.

	select {
	case ch <- resp:
	default:
		fmt.Println(len(ch))
		panic(fmt.Errorf("response channel for message ID %s is full", resp.ID))
	}
}

func (r *responseNotifier) WaitForResponse(ctx context.Context, messageId string) (*funcie.Response, error) {
	ch := r.getResponseChannel(ctx, messageId)

	select {
	case resp := <-ch:
		r.closeResponseChannel(ctx, messageId)
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *responseNotifier) getResponseChannel(ctx context.Context, messageId string) chan *funcie.Response {
	r.mapLock.Lock()
	defer r.mapLock.Unlock()

	if res, ok := r.responses[messageId]; ok {
		return res
	}

	// While WaitForResponse will usually create the channel, it's possible for the response to be received before
	// WaitForResponse is called. In this case, we need to create the channel here.

	ch := make(chan *funcie.Response, 1)
	r.responses[messageId] = ch

	return ch
}

func (r *responseNotifier) closeResponseChannel(ctx context.Context, messageId string) {
	r.mapLock.Lock()
	defer r.mapLock.Unlock()

	if res, ok := r.responses[messageId]; ok {
		close(res)
		r.responses[messageId] = nil
	} else {
		panic(fmt.Errorf("no response channel for message ID %s", messageId))
	}
}
