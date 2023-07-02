package bastion_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/mock"
	"reflect"
)

func MessageComparer(message *funcie.Message) interface{} {
	return mock.MatchedBy(func(m *funcie.Message) bool {
		cp := *m
		cp.ID = message.ID
		cp.Created = message.Created
		return reflect.DeepEqual(&cp, message)
	})
}
