package bastion_test

import (
	. "github.com/Kapps/funcie/pkg/bastion"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewConfigFromEnvironment(t *testing.T) {
	t.Run("should return a config with the correct values", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "password")
		t.Setenv("FUNCIE_REQUEST_TTL", "30m")

		config := NewConfigFromEnvironment()
		require.Equal(t, "localhost:6379", config.RedisAddress)
		require.Equal(t, "password", config.ListenAddress)
		require.Equal(t, 30*time.Minute, config.RequestTtl)
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
