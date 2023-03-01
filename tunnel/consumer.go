package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
	"time"
)

var ErrPubSubChannelClosed = fmt.Errorf("pubsub channel closed")

// Response represents a response to a message sent to a tunnel.
type Response struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Data is the actual message payload.
	Data []byte `json:"data"`
	// Received is the time the response was received.
	Received time.Time `json:"received"`
}

// NewResponse creates a new response with the given data and the current time as the received time.
func NewResponse(id string, data []byte) *Response {
	return &Response{
		ID:       id,
		Data:     data,
		Received: time.Now().Truncate(time.Millisecond),
	}
}

// Handler is a function that handles a message from a tunnel.
type Handler func(ctx context.Context, message *Message) (*Response, error)

// Consumer represents a consumer of a synchronous tunnel that can be used to receive messages from a publisher and send a response.
type Consumer interface {
	// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
	Consume(ctx context.Context, handler Handler) error
}

type RedisPubSub interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Close() error
}

// RedisConsumeClient is the interface that wraps the redis client methods used by the consumer.
type RedisConsumeClient interface {
	Subscribe(ctx context.Context, channels ...string) RedisPubSub
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

type redisConsumeClient struct {
	*redis.Client
}

func (c *redisConsumeClient) Subscribe(ctx context.Context, channels ...string) RedisPubSub {
	return c.Client.Subscribe(ctx, channels...)
}

// RedisConsumer represents a consumer that consumes messages from a Redis channel.
type RedisConsumer struct {
	redisClient RedisConsumeClient
	channelName string
}

// NewRedisConsumer creates a new RedisConsumer that consumes messages from the given channel.
// This implementation takes in a *redis.Client instead of a RedisConsumeClient so that it can be used with a real Redis client.
func NewRedisConsumer(redisClient *redis.Client, channelName string) Consumer {
	wrappedRedis := &redisConsumeClient{
		Client: redisClient,
	}
	return &RedisConsumer{
		redisClient: wrappedRedis,
		channelName: channelName,
	}
}

// NewRedisConsumerWithClient creates a new RedisConsumer that consumes messages from the given channel.
func NewRedisConsumerWithClient(redisClient RedisConsumeClient, channelName string) Consumer {
	return &RedisConsumer{
		redisClient: redisClient,
		channelName: channelName,
	}
}

// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
func (c *RedisConsumer) Consume(ctx context.Context, handler Handler) error {
	pubSub := c.redisClient.Subscribe(ctx, c.channelName)
	defer CloseOrLog(fmt.Sprintf("pubsub channel %v", c.channelName), pubSub)

	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled", "err", ctx.Err())
			return ctx.Err()
		case msg, ok := <-pubSub.Channel():
			if !ok {
				slog.Debug("pubsub channel closed")
				return ErrPubSubChannelClosed
			}

			message, err := parseMessage(msg.Payload)
			if err != nil {
				return fmt.Errorf("error parsing message: %w", err)
			}

			response, err := handler(ctx, message)
			if err != nil {
				return fmt.Errorf("error handling message: %w", err)
			}

			responseKey := GetResponseKeyForMessage(message.ID)
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

func parseMessage(message string) (*Message, error) {
	var msg Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func formatResponse(response *Response) (string, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
