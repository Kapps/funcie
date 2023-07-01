package bastion

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

// ApplicationClient allows communication between the client bastion and the client application.
type ApplicationClient interface {
	// ProcessRequest sends a request from the server bastion to the client.
	ProcessRequest(ctx context.Context, application funcie.Application, request *messages.Message) (*funcie.Response, error)
}

type httpApplicationClient struct {
	protocol string
	client   *http.Client
}

// NewHTTPApplicationClient creates a new ApplicationClient that uses the given protocol to communicate with the client application.
func NewHTTPApplicationClient(protocol string, client *http.Client) ApplicationClient {
	if protocol != "http" && protocol != "https" {
		panic(fmt.Errorf("invalid protocol %v", protocol))
	}

	return &httpApplicationClient{
		protocol: protocol,
		client:   client,
	}
}

func (h *httpApplicationClient) ProcessRequest(ctx context.Context, application funcie.Application, request *messages.Message) (*funcie.Response, error) {
	url := makeUrl(h.protocol, application.Endpoint, "process")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.Data))
	if err != nil {
		return nil, fmt.Errorf("create request to %v: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	slog.InfoCtx(ctx, "sending request to client application",
		"id", request.ID, "kind", request.Kind, "application", application.Name, "url", url)

	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request to %v: %w", url, err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	responsePayload, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %v: %w", url, err)
	}

	// TODO: Forwarding an application error when processing the request.
	resp := funcie.NewResponse(request.ID, responsePayload, nil)
	return resp, nil
}

func makeUrl(protocol string, endpoint funcie.Endpoint, path string) string {
	return fmt.Sprintf("%v://%v/%v", protocol, endpoint, path)
}
