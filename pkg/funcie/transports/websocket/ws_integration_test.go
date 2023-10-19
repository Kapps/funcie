package websocket_test

import (
	"github.com/go-faker/faker/v4"
	"testing"
)

func TestIntegration(t *testing.T) {
	/*authToken := newAuthToken()

	server := NewServer(WithBearerAuthorizationHandler(authToken))
	go func() {

	}*/
}

func newAuthToken() string {
	res := faker.Jwt()
	return res
}
