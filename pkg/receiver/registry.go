package receiver

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"log/slog"
	"sync"
)

type memoryApplicationRegistry struct {
	//registeredApplications map[string]*funcie.Application
	registeredApplications sync.Map
}

// NewMemoryApplicationRegistry creates a new in-memory application registry.
func NewMemoryApplicationRegistry() funcie.ApplicationRegistry {
	return &memoryApplicationRegistry{}
}

func (r *memoryApplicationRegistry) Register(ctx context.Context, application *funcie.Application) error {
	_, exists := r.registeredApplications.Swap(application.Name, application)
	if exists {
		slog.WarnContext(ctx,
			"application already registered; overwriting",
			"application", application.Name, "previous", application.Endpoint,
			"new", application.Endpoint,
		)
	}
	return nil
}

func (r *memoryApplicationRegistry) Unregister(_ context.Context, applicationName string) error {
	_, exists := r.registeredApplications.LoadAndDelete(applicationName)
	if !exists {
		return fmt.Errorf("application %s not registered", applicationName)
	}

	return nil
}

func (r *memoryApplicationRegistry) GetApplication(_ context.Context, applicationName string) (*funcie.Application, error) {
	application, ok := r.registeredApplications.Load(applicationName)
	if !ok {
		return nil, funcie.ErrApplicationNotFound
	}
	return application.(*funcie.Application), nil
}
