package bastion_test

import (
	"github.com/Kapps/funcie/cmd/client-bastion/bastion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfigFromEnvironment(t *testing.T) {
	t.Run("with all environment variables set", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "redis://localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "localhost:8080")
		t.Setenv("FUNCIE_BASE_CHANNEL_NAME", "override")

		config := bastion.NewConfigFromEnvironment()

		assert.Equal(t, "redis://localhost:6379", config.RedisAddress)
		assert.Equal(t, "localhost:8080", config.ListenAddress)
		assert.Equal(t, "override", config.BaseChannelName)
	})

	t.Run("with only required environment variables set", func(t *testing.T) {
		t.Setenv("FUNCIE_REDIS_ADDRESS", "redis://localhost:6379")
		t.Setenv("FUNCIE_LISTEN_ADDRESS", "localhost:8080")
		t.Setenv("FUNCIE_BASE_CHANNEL_NAME", "")

		config := bastion.NewConfigFromEnvironment()

		assert.Equal(t, "redis://localhost:6379", config.RedisAddress)
		assert.Equal(t, "localhost:8080", config.ListenAddress)
		assert.Equal(t, "funcie:requests", config.BaseChannelName)
	})

	t.Run("with no environment variables set", func(t *testing.T) {
		assert.Panics(t, func() {
			bastion.NewConfigFromEnvironment()
		})
	})
}
