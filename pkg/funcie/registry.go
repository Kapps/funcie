package funcie

import (
	"context"
	"errors"
)

// ErrApplicationNotFound is returned when an application is not found.
var ErrApplicationNotFound = errors.New("application not found")

// Application represents a registered application that can have requests routed to it.
type Application struct {
	// Name is the name of the application.
	Name string `json:"name"`
	// Endpoint is the address to send requests to.
	Endpoint string `json:"endpoint"`
}

// NewApplication creates a new Application with the given name and endpoint.
func NewApplication(name string, endpoint string) *Application {
	return &Application{
		Name:     name,
		Endpoint: endpoint,
	}
}

// ApplicationRegistry is a service that can register and unregister applications.
type ApplicationRegistry interface {
	// Register registers the given application.
	Register(ctx context.Context, application *Application) error
	// Unregister unregisters the application with the given name.
	Unregister(ctx context.Context, applicationName string) error
	// GetApplication gets the application with the given name.
	// If no application is registered with the given name, ErrApplicationNotFound is returned.
	GetApplication(ctx context.Context, applicationName string) (*Application, error)
}
