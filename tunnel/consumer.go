package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
	"time"
)

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
	Channel() <-chan *redis.Message
	Close() error
}

// RedisConsumeClient is the interface that wraps the redis client methods used by the consumer.
type RedisConsumeClient interface {
	Subscribe(ctx context.Context, channels ...string) RedisPubSub
	BRPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

// RedisConsumer represents a consumer that consumes messages from a Redis channel.
type RedisConsumer struct {
	redisClient RedisConsumeClient
	channelName string
}

// NewRedisConsumer creates a new RedisConsumer that consumes messages from the given channel.
func NewRedisConsumer(redisClient RedisConsumeClient, channelName string) Consumer {
	return &RedisConsumer{
		redisClient: redisClient,
		channelName: channelName,
	}
}

// Consume consumes a message from the tunnel, processes it, and sends the response to the other side.
func (c *RedisConsumer) Consume(ctx context.Context, handler Handler) error {
	pubSub := c.redisClient.Subscribe(ctx, c.channelName)
	defer func() {
		err := pubSub.Close()
		if err != nil {
			slog.Error("error closing pubsub", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-pubSub.Channel():
			if !ok {
				return nil
			}

			message, err := parseMessage(msg.Payload)
			if err != nil {
				return err
			}
			response, err := handler(ctx, message)
			if err != nil {
				return err
			}

			responseKey := GetResponseKeyForMessage(message.ID)
			responseData, err := formatResponse(response)
			if err != nil {
				return fmt.Errorf("error formatting response: %w", err)
			}

			cmd := c.redisClient.BRPush(ctx, responseKey, responseData)
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
