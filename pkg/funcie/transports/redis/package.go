// Package redis provides a Redis, queue-based, transport for funcie.
// Consumers subscribe to a channel that is named after the application that they are consuming for.
// The publishers send messages to that channel, which returns whether there are any active consumers.
// Consumers wait for messages on the channel, and then process them synchronously.
// This means consumers have a blocking, ordered, queue and can process only a single request at a time.
// Producers on the other hand, can send many messages at once, and will receive a response for each one.
package redis
