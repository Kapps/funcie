package bastion_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"reflect"
	"testing"
)

func MessageComparer(t *testing.T, message *funcie.Message) interface{} {
	return func(m *funcie.Message) bool {
		cp := *m
		cp.ID = message.ID
		cp.Created = message.Created
		return reflect.DeepEqual(&cp, message)
	}
}
