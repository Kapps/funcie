package redis

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type PubSub interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Close() error
}

// ConsumeClient is the interface that wraps the redis client methods used by the consumer.
type ConsumeClient interface {
	Subscribe(ctx context.Context, channels ...string) PubSub
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

type redisConsumeClient struct {
	*redis.Client
}

func (c *redisConsumeClient) Subscribe(ctx context.Context, channels ...string) PubSub {
	return c.Client.Subscribe(ctx, channels...)
}

// Consumer represents a consumer that consumes messages from a Redis channel.
type Consumer struct {
	redisClient ConsumeClient
	channelName string
}

// NewConsumer creates a new RedisConsumer that consumes messages from the given channel.
// This implementation takes in a *redis.Client instead of a RedisConsumeClient so that it can be used with a real Redis client.
func NewConsumer(redisClient *redis.Client, channelName string) funcie.Consumer {
	wrappedRedis := &redisConsumeClient{
		Client: redisClient,
	}
	return &Consumer{
		redisClient: wrappedRedis,
		channelName: channelName,
	}
}

// NewConsumerWithClient creates a new RedisConsumer that consumes messages from the given channel.
func NewConsumerWithClient(redisClient ConsumeClient, channelName string) funcie.Consumer {
	return &Consumer{
		redisClient: redisClient,
		channelName: channelName,
	}
}

// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
func (c *Consumer) Consume(ctx context.Context, handler funcie.Handler) error {
	pubSub := c.redisClient.Subscribe(ctx, c.channelName)
	defer funcie.CloseOrLog(fmt.Sprintf("pubsub channel %v", c.channelName), pubSub)

	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled", "err", ctx.Err())
			return ctx.Err()
		case msg, ok := <-pubSub.Channel():
			if !ok {
				slog.Debug("pubsub channel closed")
				return funcie.ErrPubSubChannelClosed
			}

			message, err := parseMessage(msg.Payload)
			if err != nil {
				return fmt.Errorf("error parsing message: %w", err)
			}

			response, err := handler(ctx, message)
			if err != nil {
				return fmt.Errorf("error handling message: %w", err)
			}

			responseKey := funcie.GetResponseKeyForMessage(message.ID)
			responseData, err := formatResponse(response)
			if err != nil {
				return fmt.Errorf("error formatting response: %w", err)
			}

			cmd := c.redisClient.RPush(ctx, responseKey, responseData)
			if err := cmd.Err(); err != nil {
				return fmt.Errorf("error pushing response to queue: %w", err)
			}
		}
	}
}

func parseMessage(message string) (*funcie.Message, error) {
	var msg funcie.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func formatResponse(response *funcie.Response) (string, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
*/
