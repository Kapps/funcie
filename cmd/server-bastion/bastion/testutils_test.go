package bastion_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func MessageComparer(message *funcie.Message) interface{} {
	return mock.MatchedBy(func(m *funcie.Message) bool {
		cp := *m
		cp.ID = message.ID
		cp.Created = message.Created
		return reflect.DeepEqual(&cp, message)
	})
}

func ResponseComparer(response *funcie.Response) interface{} {
	return mock.MatchedBy(func(r *funcie.Response) bool {
		cp := *r
		cp.ID = response.ID
		cp.Received = response.Received
		return reflect.DeepEqual(&cp, response)
	})
}

func RequireEqualResponse[T any](t *testing.T, expected, actual *funcie.ResponseBase[T]) {
	t.Helper()
	cp := *expected
	cp.ID = actual.ID
	cp.Received = actual.Received
	require.Equal(t, &cp, actual)
}
