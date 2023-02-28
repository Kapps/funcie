package tunnel

import (
	"encoding/json"
	"fmt"
)

// GetResponseKeyForMessage returns the Redis key for the response of a message
func GetResponseKeyForMessage(messageId string) string {
	if messageId == "" {
		panic("messageId cannot be empty")
	}
	return fmt.Sprintf("funcie:resp:%v", messageId)
}

// MustSerialize serializes the given value to JSON, or panics if it fails.
func MustSerialize(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
