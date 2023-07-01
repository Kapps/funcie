package funcie

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// Endpoint represents a target destination for funcie requests or bastions.
type Endpoint struct {
	// Scheme is the protocol scheme of the endpoint.
	// This is usually "http" or "https".
	Scheme string `json:"protocol"`
	// Host is the host name or IP address of the endpoint.
	// This does not include the protocol or port.
	Host string `json:"host"`
	// Port is the port number of the endpoint.
	Port int `json:"port"`
}

// NewEndpoint creates a new Endpoint with the given scheme, host, and port.
func NewEndpoint(scheme string, host string, port int) Endpoint {
	return Endpoint{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}
}

// String returns a URL representation of the endpoint.
func (e Endpoint) String() string {
	return fmt.Sprintf("%v://%v:%v", e.Scheme, e.Host, e.Port)
}

// NewEndpointFromAddress creates a new funcie Endpoint from parsing the given address.
// Example: https://127.0.0.1:8080
func NewEndpointFromAddress(address string) (Endpoint, error) {
	parsed, err := url.Parse(address)
	if err != nil {
		return Endpoint{}, fmt.Errorf("parsing address %v: %w", address, err)
	}

	host, port, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		return Endpoint{}, fmt.Errorf("splitting host and port from address %v: %w", address, err)
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return Endpoint{}, fmt.Errorf("converting port %v to int: %w", port, err)
	}

	return NewEndpoint(parsed.Scheme, host, parsedPort), nil
}

// MustNewEndpointFromAddress creates a new funcie Endpoint from parsing the given address.
func MustNewEndpointFromAddress(address string) Endpoint {
	endpoint, err := NewEndpointFromAddress(address)
	if err != nil {
		panic(err)
	}

	return endpoint
}
