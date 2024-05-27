package funcli

import (
	"context"
	"fmt"
	"net/http"
	ws "nhooyr.io/websocket"

	"io"
	"net"
)

// TunnelOptions provides optional arguments for creating a network tunnel.
type TunnelOptions struct {
	// Headers is a map of headers to send with the connection request.
	Headers http.Header
}

// Tunneller is an interface for creating a tunnel to a remote host.
type Tunneller interface {
	// OpenTunnel starts a tunnel to a remote host on the given port locally.
	OpenTunnel(ctx context.Context, endpoint string, localPort int, opts *TunnelOptions) error
}

type webhookTunnel struct {
}

// NewWebhookTunneller creates a new WebhookTunnel.
func NewWebhookTunneller() Tunneller {
	return &webhookTunnel{}
}

func (t *webhookTunnel) OpenTunnel(ctx context.Context, endpoint string, localPort int, opts *TunnelOptions) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", localPort))
	if err != nil {
		return fmt.Errorf("failed to dial local connection: %w", err)
	}
	defer func() { _ = listener.Close() }()

	for {
		fmt.Println("Waiting for connection on port", localPort)
		localConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: failed to accept local connection", "error", err)
			continue
		}

		fmt.Println("Accepted connection on port", localPort)

		wsConn, _, err := ws.Dial(ctx, endpoint, &ws.DialOptions{
			HTTPHeader: opts.Headers,
			//Subprotocols: []string{"binary"},
			CompressionMode: ws.CompressionDisabled,
		})
		if err != nil {
			return fmt.Errorf("failed to dial WebSocket connection: %w", err)
		}

		fmt.Println("Connected to WebSocket")

		go func() {
			err := t.forward(ctx, localConn, wsConn)
			if err != nil {
				fmt.Println("Error: failed to forward connection", "error", err)
			}
			fmt.Println("Done forwarding connection")
		}()
	}
}

func (t *webhookTunnel) forward(ctx context.Context, localConn net.Conn, wsConn *ws.Conn) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	//defer cancel()
	//defer func() { _ = localConn.Close() }()
	//defer func() { _ = wsConn.Close(ws.StatusGoingAway, "exiting") }()

	errChan := make(chan error, 1)

	go func() {
		for {
			fmt.Println("Reading from WebSocket")
			_, message, err := wsConn.Read(cancelCtx)
			if err != nil {
				fmt.Println("Error reading from WebSocket: %w", err)
				errChan <- fmt.Errorf("failed to read from WebSocket connection: %w", err)
				cancel()
				return
			}

			fmt.Println("Writing to local connection")
			_, err = localConn.Write(message)
			if err != nil {
				fmt.Println("Error writing to local connection: %w", err)
				errChan <- fmt.Errorf("failed to write to local connection: %w", err)
				cancel()
				return
			}
		}
	}()

	buffer := make([]byte, 20)
	for {
		if len(errChan) > 0 {
			fmt.Println("Error channel has data")
			return <-errChan
		}

		fmt.Println("Reading from local connection")
		n, err := localConn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from local connection: %w", err)
				return fmt.Errorf("failed to read from local connection: %w", err)
			}
			fmt.Println("EOF from local connection")
			return nil
		}

		fmt.Println("Read", n, "bytes from local connection")

		fmt.Println("Writing to WebSocket")
		err = wsConn.Write(ctx, ws.MessageBinary, buffer[:n])
		if err != nil {
			fmt.Println("Error writing to WebSocket: %w", err)
			return fmt.Errorf("failed to write to WebSocket connection: %w", err)
		}

		fmt.Println("Wrote", n, "bytes to WebSocket")
	}
}
