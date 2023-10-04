package redis

import (
	"errors"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
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

// IsNoHandlerFound returns true if the given error is a ErrNoHandlerFound,
// or if the given response is a NoHandlerFound response.
func IsNoHandlerFound(err error, resp *funcie.Response) bool {
	// Gross.
	return errors.Is(err, utils.ErrNoHandlerFound) ||
		(err == nil && resp.Error != nil && (resp.Error.Message == utils.ErrNoHandlerFound.Error() || resp.Error.Message == funcie.ErrNoActiveConsumer.Error()))
}
