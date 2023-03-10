package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/redis/go-redis/v9"
	"time"
)

type PublishClient interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
}

type Publisher struct {
	redisClient PublishClient
	channelName string
}

// NewPublisher creates a new RedisPublisher that publishes messages to the given channel.
func NewPublisher(redisClient PublishClient, channelName string) *Publisher {
	return &Publisher{redisClient: redisClient, channelName: channelName}
}

func (p *Publisher) Publish(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	// Publish the message to the channel.
	messageContents, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	pub := p.redisClient.Publish(ctx, p.channelName, messageContents)
	if err := pub.Err(); err != nil {
		return nil, fmt.Errorf("failed to publish message to channel %s: %w", p.channelName, err)
	}

	consumers, err := pub.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get result of publish: %w", err)
	}

	if consumers == 0 {
		return nil, funcie.ErrNoActiveConsumer
	}

	// Wait for a response from the consumer.
	responseKey := funcie.GetResponseKeyForMessage(message.ID)
	resp, err := p.redisClient.BRPop(ctx, message.Ttl, responseKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get response within ttl: %w", err)
	}

	// The response should be a list of two items, the first being the key and the second being the value.
	if len(resp) == 0 {
		return nil, fmt.Errorf("no response received within ttl of %s", message.Ttl)
	}

	if len(resp) != 2 {
		panic(fmt.Sprintf("expected response to be a list of two items, got %d", len(resp)))
	}

	// First entry is key; value should be the serialized response data.
	data := []byte(resp[1])

	var response *funcie.Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
