package receiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"golang.org/x/exp/slog"
	"net/http"
)

type applicationRegistrar struct {
	registeredApplications map[string]*funcie.Application
	bastionAddress         string
}

// NewBastionBackedMemoryApplicationRegistry creates a new ApplicationRegistry that stores locally in memory and
// forwards the changes to a bastion server.
func NewBastionBackedMemoryApplicationRegistry(bastionAddress string) funcie.ApplicationRegistry {
	return &applicationRegistrar{
		registeredApplications: make(map[string]*funcie.Application),
		bastionAddress:         bastionAddress,
	}
}

func (r *applicationRegistrar) Register(ctx context.Context, application *funcie.Application) error {
	_, ok := r.registeredApplications[application.Name]
	if ok {
		slog.Warn(
			"application %s already registered; overwriting",
			"application", application.Name, "previous", application.Endpoint,
			"new", application.Endpoint,
		)
	}
	r.registeredApplications[application.Name] = application
	return nil
}

func (r *applicationRegistrar) Unregister(ctx context.Context, applicationName string) error {
	_, ok := r.registeredApplications[applicationName]
	if !ok {
		return fmt.Errorf("application %s not registered", applicationName)
	}

	delete(r.registeredApplications, applicationName)
	return nil
}

func (r *applicationRegistrar) GetApplication(ctx context.Context, applicationName string) (*funcie.Application, error) {
	application, ok := r.registeredApplications[applicationName]
	if !ok {
		return nil, funcie.ErrApplicationNotFound
	}
	return application, nil
}

func post[Response any](ctx context.Context, url string, payload interface{}) (*Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer funcie.CloseOrLog("response body", resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
