package funcie_test

import (
	. "github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetResponseKeyForMessage(t *testing.T) {
	t.Parallel()

	t.Run("should return a key for a message", func(t *testing.T) {
		t.Parallel()

		messageId := "message-id"
		expectedKey := "funcie:resp:message-id"

		key := GetResponseKeyForMessage(messageId)

		require.Equal(t, expectedKey, key)
	})

	t.Run("should panic if the message ID is empty", func(t *testing.T) {
		t.Parallel()

		messageId := ""

		require.Panics(t, func() {
			GetResponseKeyForMessage(messageId)
		})
	})
}
