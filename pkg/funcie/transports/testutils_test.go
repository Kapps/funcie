package transports_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func MessageComparer[T any](message *funcie.MessageBase[T]) interface{} {
	return mock.MatchedBy(func(m *funcie.MessageBase[T]) bool {
		cp := *m
		cp.ID = message.ID
		cp.Created = message.Created
		return reflect.DeepEqual(&cp, message)
	})
}

func RequireEqualResponse[T any](t *testing.T, expected *funcie.ResponseBase[T], actual *funcie.ResponseBase[T]) {
	t.Helper()
	diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(funcie.ResponseBase[T]{}, "Received"))
	if diff != "" {
		t.Errorf("unexpected response (-want +got):\n%s", diff)
	}
}
