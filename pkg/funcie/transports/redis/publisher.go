package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
	"time"
)

var ttl = time.Minute * 5

type PublishClient interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
}

type redisPublisher struct {
	redisClient     PublishClient
	baseChannelName string
}

// NewPublisher creates a new RedisPublisher that publishes messages to the given channel.
func NewPublisher(redisClient PublishClient, baseChannelName string) funcie.Publisher {
	return &redisPublisher{
		redisClient:     redisClient,
		baseChannelName: baseChannelName,
	}
}

func (p *redisPublisher) Publish(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	channelName := GetChannelNameForApplication(p.baseChannelName, message.Application)

	messageContents, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	slog.InfoCtx(ctx, "publishing message to channel", "channel", channelName, "message", message.ID)

	pub := p.redisClient.Publish(ctx, channelName, messageContents)
	if err := pub.Err(); err != nil {
		return nil, fmt.Errorf("failed to publish message to channel %s: %w", message.Application, err)
	}

	slog.DebugCtx(ctx, "published message to channel", "channel", channelName, "message", message.ID)

	consumers, err := pub.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get result of publish: %w", err)
	}

	slog.DebugCtx(ctx, "received publish result", "consumers", consumers)

	if consumers == 0 {
		return nil, funcie.ErrNoActiveConsumer
	}

	// Wait for a response from the consumer.
	responseKey := GetResponseKeyForMessage(p.baseChannelName, message.ID)
	resp, err := p.redisClient.BRPop(ctx, ttl, responseKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get response from consumer: %w", err)
	}

	if len(resp) == 0 {
		return nil, funcie.ErrNoActiveConsumer
	}

	if len(resp) != 2 {
		panic(fmt.Sprintf("expected response to be a list of two items, got %d", len(resp)))
	}

	// First entry is key; value should be the serialized response data.
	responseContents := []byte(resp[1])

	var response funcie.Response
	if err := json.Unmarshal(responseContents, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from consumer: %w", err)
	}

	return &response, nil
}
