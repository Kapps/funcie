package bastion_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/mock"
	"reflect"
)

func MessageComparer[T any](message *funcie.MessageBase[T]) interface{} {
	return mock.MatchedBy(func(m *funcie.MessageBase[T]) bool {
		cp := *m
		cp.ID = message.ID
		cp.Created = message.Created
		return reflect.DeepEqual(&cp, message)
	})
}
