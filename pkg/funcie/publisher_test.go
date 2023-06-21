package funcie_test

import (
	. "github.com/Kapps/funcie/pkg/funcie/messages"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	t.Parallel()

	t.Run("should return a new message with a unique ID", func(t *testing.T) {
		t.Parallel()

		m1 := NewMessage("foo", MessageKindDispatch, []byte("hello"), time.Second)
		m2 := NewMessage("foo", MessageKindDispatch, []byte("hello"), time.Second)

		require.NotEqual(t, m1.ID, m2.ID)
	})

	t.Run("should return a new message with the given data", func(t *testing.T) {
		t.Parallel()

		m := NewMessage("foo", MessageKindDispatch, []byte("hello"), time.Second)

		require.Equal(t, m.Data, []byte("hello"))
	})
}
