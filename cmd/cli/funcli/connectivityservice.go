package funcli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ConnectivityService provides utilities for testing / awaiting internet connectivity.
type ConnectivityService interface {
	// WaitForConnectivity waits for the given endpoint to be reachable, or the context to be done.
	WaitForConnectivity(ctx context.Context, endpoint string) error
}

type HttpConnectivityServiceOptions struct {
	RetryInterval time.Duration
}

type httpConnectivityService struct {
	opts HttpConnectivityServiceOptions
}

type HttpConnectivityServiceOptionSetter func(*HttpConnectivityServiceOptions)

// WithRetryInterval sets the retry interval for the HttpConnectivityService.
func WithRetryInterval(interval time.Duration) HttpConnectivityServiceOptionSetter {
	return func(opts *HttpConnectivityServiceOptions) {
		opts.RetryInterval = interval
	}
}

// NewHttpConnectivityService creates a new HttpConnectivityService with optional settings.
func NewHttpConnectivityService(opts ...HttpConnectivityServiceOptionSetter) ConnectivityService {
	config := &HttpConnectivityServiceOptions{
		RetryInterval: 1 * time.Second,
	}

	for _, setter := range opts {
		setter(config)
	}

	return &httpConnectivityService{
		opts: *config,
	}
}

func (s *httpConnectivityService) WaitForConnectivity(ctx context.Context, endpoint string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := http.NewRequestWithContext(ctx, http.MethodOptions, endpoint, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			req.Close = true

			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				_ = resp.Body.Close()
				return nil
			}
			if errors.Is(err, http.ErrServerClosed) || errors.Is(err, http.ErrHandlerTimeout) {
				time.Sleep(s.opts.RetryInterval)
				continue
			}
			return fmt.Errorf("failed to connect to %s: %w", endpoint, err)
		}
	}
}
