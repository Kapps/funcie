package funcie

import (
	"fmt"
	"net"
	"strconv"
)

// Endpoint represents a target destination for funcie requests or bastions.
type Endpoint struct {
	// Host is the host name or IP address of the endpoint.
	Host string `json:"host"`
	// Port is the port number of the endpoint.
	Port int `json:"port"`
}

// NewEndpoint creates a new Endpoint with the given host and port.
func NewEndpoint(host string, port int) Endpoint {
	return Endpoint{
		Host: host,
		Port: port,
	}
}

// NewEndpointFromAddress creates a new funcie Endpoint from parsing the given address.
// Example: 127.0.0.1:8080
func NewEndpointFromAddress(address string) (Endpoint, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return Endpoint{}, fmt.Errorf("splitting address %v: %w", address, err)
	}

	numericPort, err := strconv.Atoi(port)
	if err != nil {
		return Endpoint{}, fmt.Errorf("parsing port %v: %w", port, err)
	}

	return NewEndpoint(host, numericPort), nil
}
