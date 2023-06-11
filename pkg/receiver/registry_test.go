package receiver_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	. "github.com/Kapps/funcie/pkg/receiver"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryApplicationRegistry_Integration(t *testing.T) {
	registry := NewMemoryApplicationRegistry()
	app := &funcie.Application{
		Name:     "test",
		Endpoint: "ep",
	}

	t.Run("loading with no applications", func(t *testing.T) {
		app, err := registry.GetApplication(nil, "test")
		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
		require.Nil(t, app)
	})

	t.Run("registering an application", func(t *testing.T) {
		err := registry.Register(nil, app)
		require.NoError(t, err)
	})

	t.Run("loading a registered application", func(t *testing.T) {
		loaded, err := registry.GetApplication(nil, "test")
		require.NoError(t, err)
		require.Equal(t, loaded, app)
	})

	t.Run("unregistering an application", func(t *testing.T) {
		err := registry.Unregister(nil, "test")
		require.NoError(t, err)
	})

	t.Run("loading an unregistered application", func(t *testing.T) {
		loaded, err := registry.GetApplication(nil, "test")
		require.ErrorIs(t, err, funcie.ErrApplicationNotFound)
		require.Nil(t, loaded)
	})
}
