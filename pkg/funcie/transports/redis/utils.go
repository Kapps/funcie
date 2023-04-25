package redis

import (
	"fmt"
)

// GetResponseKeyForMessage returns the Redis key for the response of a message
func GetResponseKeyForMessage(baseChannelName string, messageId string) string {
	if messageId == "" {
		panic("messageId cannot be empty")
	}
	return fmt.Sprintf("%v:resp:%v", baseChannelName, messageId)
}

// GetChannelNameForApplication returns the Redis channel name for the given application ID.
func GetChannelNameForApplication(baseChannelName string, applicationId string) string {
	if applicationId == "" {
		panic("applicationId cannot be empty")
	}
	return fmt.Sprintf("%v:app:%v", baseChannelName, applicationId)
}
