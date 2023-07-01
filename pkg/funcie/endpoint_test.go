package funcie_test

import (
	. "github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEndpoint(t *testing.T) {
	t.Parallel()

	endpoint := NewEndpoint("https", "host", 1234)
	require.Equal(t, "host", endpoint.Host)
	require.Equal(t, 1234, endpoint.Port)
	require.Equal(t, "https", endpoint.Scheme)

	require.Equal(t, "https://host:1234", endpoint.String())
}

func TestNewEndpointFromAddress(t *testing.T) {
	t.Parallel()

	t.Run("should return an endpoint with the correct values", func(t *testing.T) {
		endpoint, err := NewEndpointFromAddress("https://127.0.0.1:8080")
		require.NoError(t, err)

		require.Equal(t, "https", endpoint.Scheme)
		require.Equal(t, "127.0.0.1", endpoint.Host)
		require.Equal(t, 8080, endpoint.Port)
	})

	t.Run("should return an error if the address is invalid", func(t *testing.T) {
		_, err := NewEndpointFromAddress("invalid")
		require.Error(t, err)
	})

	t.Run("should return an error if the port is invalid", func(t *testing.T) {
		_, err := NewEndpointFromAddress("http://127.0.0.1:invalid")
		require.Error(t, err)
	})

	t.Run("should return an error if the port is missing", func(t *testing.T) {
		_, err := NewEndpointFromAddress("http://127.0.0.1")
		require.Error(t, err)
	})

	t.Run("should return an error if the scheme is missing", func(t *testing.T) {
		_, err := NewEndpointFromAddress("127.0.0.1:8080")
		require.Error(t, err)
	})
}
