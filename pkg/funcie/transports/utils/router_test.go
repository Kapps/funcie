package utils

import (
	"context"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandlerRouter_TestHandler(t *testing.T) {
	t.Parallel()

	id1, id2 := faker.Word(), faker.Word()
	application := faker.Word()

	called := false
	handler := funcie.Handler(func(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
		assert.Equalf(t, id2, message.ID, "message id was not passed to handler")

		called = true

		return &funcie.Response{
			ID:       id1,
			Data:     nil,
			Error:    nil,
			Received: time.Time{},
		}, nil
	})
	router := NewClientHandlerRouter()

	err := router.AddClientHandler(application, handler)
	assert.NoError(t, err)

	response, err := router.Handle(context.Background(), &funcie.Message{
		ID:          id2,
		Application: application,
	})
	assert.NoError(t, err)
	assert.Equalf(t, id1, response.ID, "response id was not returned from handler")

	err = router.RemoveClientHandler(application)
	assert.NoError(t, err)

	response, err = router.Handle(context.Background(), &funcie.Message{
		ID:          id2,
		Application: application,
	})
	assert.Error(t, err)

	assert.Truef(t, called, "handler was not called")
}
