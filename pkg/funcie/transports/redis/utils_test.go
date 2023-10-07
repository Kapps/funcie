package redis_test

import (
	. "github.com/Kapps/funcie/pkg/funcie/transports/redis"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetResponseKeyForMessage(t *testing.T) {
	t.Parallel()

	baseChannelName := "funcie"

	t.Run("should return a key for a message", func(t *testing.T) {
		t.Parallel()

		messageId := "message-id"
		expectedKey := "funcie:resp:message-id"

		key := GetResponseKeyForMessage(baseChannelName, messageId)

		require.Equal(t, expectedKey, key)
	})

	t.Run("should panic if the message ID is empty", func(t *testing.T) {
		t.Parallel()

		messageId := ""

		require.Panics(t, func() {
			GetResponseKeyForMessage(baseChannelName, messageId)
		})
	})
}

func TestGetChannelNameForApplication(t *testing.T) {
	t.Parallel()

	baseChannelName := "funcie"

	t.Run("should return a channel name for an application", func(t *testing.T) {
		t.Parallel()

		applicationId := "application-id"
		expectedChannelName := "funcie:app:application-id"

		channelName := GetChannelNameForApplication(baseChannelName, applicationId)

		require.Equal(t, expectedChannelName, channelName)
	})

	t.Run("should panic if the application ID is empty", func(t *testing.T) {
		t.Parallel()

		applicationId := ""

		require.Panics(t, func() {
			GetChannelNameForApplication(baseChannelName, applicationId)
		})
	})
}
