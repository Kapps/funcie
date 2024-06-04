package funcli_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kapps/funcie/cmd/cli/funcli"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"net/http/httptest"
)

func TestWaitForConnectivity_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	service := funcli.NewHttpConnectivityService()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := service.WaitForConnectivity(ctx, server.URL)
	require.NoError(t, err)
}

func TestWaitForConnectivity_ContextDone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	service := funcli.NewHttpConnectivityService(funcli.WithRetryInterval(50 * time.Millisecond))

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	done := make(chan struct{})

	go func() {
		time.Sleep(200 * time.Millisecond)
		server.Start()
		time.Sleep(100 * time.Millisecond)
		close(done)
	}()

	url := fmt.Sprintf("http://%s", server.Listener.Addr().String())
	err := service.WaitForConnectivity(ctx, url)
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))

	<-done
}

func TestWaitForConnectivity_DelayedServerStart(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	service := funcli.NewHttpConnectivityService(funcli.WithRetryInterval(100 * time.Millisecond))

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	go func() {
		time.Sleep(500 * time.Millisecond)
		server.Start()
	}()

	url := fmt.Sprintf("http://%s", server.Listener.Addr().String())
	err := service.WaitForConnectivity(ctx, url)
	require.NoError(t, err)
}
