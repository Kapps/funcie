package websocket

import (
	"crypto/rand"
	"testing"
)

func TestIntegration(t *testing.T) {
	/*authToken := newAuthToken()

	server := NewServer(WithBasicAuthorizationHandler(authToken))
	go func() {

	}*/
}

func newAuthToken() string {
	authToken := make([]byte, 32)
	_, err := rand.Read(authToken)
	if err != nil {
		panic(err)
	}
	return string(authToken)
}
