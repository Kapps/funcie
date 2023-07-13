package funcie_test

import (
	"encoding/json"
	. "github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMessage(t *testing.T) {
	t.Parallel()

	t.Run("should return a new message with a unique ID", func(t *testing.T) {
		t.Parallel()

		m1 := NewMessage("foo", messages.MessageKindForwardRequest, []byte("hello"))
		m2 := NewMessage("foo", messages.MessageKindForwardRequest, []byte("hello"))

		require.NotEqual(t, m1.ID, m2.ID)
	})

	t.Run("should return a new message with the given data", func(t *testing.T) {
		t.Parallel()

		m := NewMessage("foo", messages.MessageKindForwardRequest, []byte("hello"))

		require.Equal(t, m.Payload, json.RawMessage("hello"))
	})
}
