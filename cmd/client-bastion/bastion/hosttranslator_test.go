package bastion

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerHostTranslator_NoTranslationsRequired(t *testing.T) {
	ctx := context.Background()
	translator := NewDockerHostTranslator()
	setCleanup(t)

	lookupHost = func(host string) ([]string, error) {
		require.Equal(t, "host.docker.internal", host)
		dnsErr := &net.DNSError{
			Err:        "no such host",
			Name:       "host.docker.internal",
			IsNotFound: true,
		}
		return nil, dnsErr
	}

	t.Run("should not require translation", func(t *testing.T) {
		req, err := translator.IsHostTranslationRequired(ctx)
		require.NoError(t, err)
		require.False(t, req)
	})

	t.Run("should not translate", func(t *testing.T) {
		result, err := translator.TranslateLocalHostToResolvedHost(ctx, "localhost")
		require.NoError(t, err)
		require.Equal(t, "localhost", result)
	})
}

func TestDockerHostTranslator_TranslationsRequired(t *testing.T) {
	ctx := context.Background()
	translator := NewDockerHostTranslator()
	setCleanup(t)

	lookupHost = func(host string) ([]string, error) {
		require.Equal(t, "host.docker.internal", host)
		return []string{"192.168.65.2"}, nil
	}

	t.Run("should require translation", func(t *testing.T) {
		req, err := translator.IsHostTranslationRequired(ctx)
		require.NoError(t, err)
		require.True(t, req)
	})

	t.Run("should translate", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"localhost", "host.docker.internal"},
			{"127.0.0.1", "host.docker.internal"},
			{"0.0.0.0", "host.docker.internal"},
			{"example.com", "example.com"},
		}

		for _, test := range tests {
			result, err := translator.TranslateLocalHostToResolvedHost(ctx, test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, result)
		}
	})
}

func setCleanup(t *testing.T) {
	originalLookupHost := lookupHost
	t.Cleanup(func() {
		lookupHost = originalLookupHost
	})
}
