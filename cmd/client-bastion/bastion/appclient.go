package bastion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

// ApplicationClient allows communication between the client bastion and the client application.
type ApplicationClient interface {
	// ProcessRequest sends a request from the server bastion to the client.
	ProcessRequest(ctx context.Context, application funcie.Application, request *funcie.Message) (*funcie.Response, error)
}

type httpApplicationClient struct {
	client *http.Client
}

// NewHTTPApplicationClient creates a new ApplicationClient that uses the given HttpClient to communicate with the client application.
func NewHTTPApplicationClient(client *http.Client) ApplicationClient {
	return &httpApplicationClient{
		client: client,
	}
}

func (h *httpApplicationClient) ProcessRequest(ctx context.Context, application funcie.Application, request *funcie.Message) (*funcie.Response, error) {
	url := makeUrl(application.Endpoint, "process")

	serialized, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("serialize request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(serialized))
	if err != nil {
		return nil, fmt.Errorf("create request to %v: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	slog.InfoCtx(ctx, "sending request to client application",
		"id", request.ID, "kind", request.Kind, "application", application.Name, "url", url)

	slog.DebugCtx(ctx, "sending message", "message", string(serialized))

	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request to %v: %w", url, err)
	}

	defer func() { _ = httpResponse.Body.Close() }()

	responsePayload, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %v: %w", url, err)
	}

	var response funcie.Response
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		slog.WarnCtx(ctx, "failed to deserialize response from client application",
			"error", err, "payload", string(responsePayload))
		return nil, fmt.Errorf("deserialize response from %v: %w", url, err)
	}

	return &response, nil
}

func makeUrl(endpoint funcie.Endpoint, path string) string {
	return fmt.Sprintf("%v/%v", endpoint.String(), path)
}
