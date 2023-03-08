package bastion_test

import (
	. "github.com/Kapps/funcie/pkg/bastion"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	config := NewConfig()
	require.NotNilf(t, config, "config should not be nil")
}

func TestNewConfigFromEnvironment(t *testing.T) {
	t.Run("should return a config with the correct values", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "password")
		t.Setenv("FUNCIE_REQUEST_TTL", "30m")
		t.Setenv("FUNCIE_REQUEST_CHANNEL", "channel")
		t.Setenv("FUNCIE_RESPONSE_KEY_PREFIX", "prefix:")

		config := NewConfigFromEnvironment()
		require.Equal(t, "localhost:6379", config.RedisAddress)
		require.Equal(t, "password", config.ListenAddress)
		require.Equal(t, 30*time.Minute, config.RequestTtl)
		require.Equal(t, "channel", config.RequestChannel)
		require.Equal(t, "prefix:", config.ResponseKeyPrefix)
	})

	t.Run("should panic if FUNCIE_REDIS_ADDRESS is not set", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "password")
		t.Setenv("FUNCIE_REQUEST_TTL", "15m")

		require.Panics(t, func() {
			NewConfigFromEnvironment()
		})
	})

	t.Run("should panic if FUNCIE_LISTEN_ADDRESS is not set", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "")
		t.Setenv("FUNCIE_REQUEST_TTL", "15m")

		require.Panics(t, func() {
			NewConfigFromEnvironment()
		})
	})

	t.Run("should default FUNCIE_REQUEST_TTL to 15 minutes if not set", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "password")
		t.Setenv("FUNCIE_REQUEST_TTL", "")

		config := NewConfigFromEnvironment()
		require.Equal(t, 15*time.Minute, config.RequestTtl)
	})
}
