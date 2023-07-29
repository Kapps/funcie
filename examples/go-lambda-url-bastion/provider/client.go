package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"io"
	"net/http"
	"net/url"
)

// BastionClient is a client that can send requests to a server bastion.
type BastionClient interface {
	// SendRequest sends a request to the bastion.
	SendRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error)
}

type httpBastionClient struct {
	client   *http.Client
	endpoint url.URL
}

// NewHTTPBastionClient creates a new BastionClient that uses HTTP to communicate with the bastion.
func NewHTTPBastionClient(endpoint url.URL) BastionClient {
	return &httpBastionClient{
		client:   &http.Client{},
		endpoint: endpoint,
	}
}

func (c *httpBastionClient) SendRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	httpResp, err := c.client.Post(c.endpoint.String(), "application/json", bytes.NewReader(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer func() { _ = httpResp.Body.Close() }()

	var response funcie.Response
	responseData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return &response, nil
}
