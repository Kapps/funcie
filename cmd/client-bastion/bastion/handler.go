package bastion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"io"
	"net/http"
)

// Handler allows the handling of incoming valid Bastion requests.
type Handler interface {
	// Register registers the given application.
	Register(ctx context.Context, application *funcie.Application) error
	// Unregister unregisters the application with the given name.
	Unregister(ctx context.Context, applicationName string) error
	// ForwardRequest forwards the given request to the application specified in the request.
	ForwardRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error)
}

type handler struct {
	registry funcie.ApplicationRegistry
}

// NewHandler creates a new Handler that can register and unregister applications and forward requests.
func NewHandler(registry funcie.ApplicationRegistry) Handler {
	return &handler{
		registry: registry,
	}
}

func (h *handler) Register(ctx context.Context, application *funcie.Application) error {
	err := h.registry.Register(ctx, application)
	if err != nil {
		return fmt.Errorf("register application %v: %w", application, err)
	}

	return nil
}

func (h *handler) Unregister(ctx context.Context, applicationName string) error {
	err := h.registry.Unregister(ctx, applicationName)
	if err != nil {
		return fmt.Errorf("unregister application %v: %w", applicationName, err)
	}

	return nil
}

func (h *handler) ForwardRequest(ctx context.Context, request *funcie.Message) (*funcie.Response, error) {
	app, err := h.registry.GetApplication(ctx, request.Application)
	if err != nil {
		return nil, fmt.Errorf("getting application %v: %w", request.Application, err)
	}

	contents, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	postResponse, err := http.Post(app.Endpoint, "application/json", bytes.NewBuffer(contents))
	if err != nil {
		return nil, fmt.Errorf("posting request to app %v endpoint %v: %w", app.Name, app.Endpoint, err)
	}

	responseData, err := io.ReadAll(postResponse.Body)

	var response funcie.Response
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling response %v: %w", responseData, err)
	}

	return &response, nil
}
