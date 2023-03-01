package tunnel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

// ErrNoActiveConsumer is returned when a consumer is not active on a tunnel.
var ErrNoActiveConsumer = errors.New("no consumer is active on this tunnel")

// Message represents a message to be sent through a tunnel.
type Message struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// Data is the actual message payload.
	Data []byte `json:"data"`
	// Created is the time the message was created.
	Created time.Time `json:"created"`
	// Ttl is the time to live for this message.
	// If the message is not processed by Created+Ttl, it should be discarded.
	Ttl time.Duration `json:"ttl"`
}

func NewMessage(data []byte, ttl time.Duration) *Message {
	return &Message{
		ID:      newId(),
		Data:    data,
		Created: time.Now().Truncate(time.Millisecond),
		Ttl:     ttl,
	}
}

func newId() string {
	return uuid.New().String()
}

// Publisher represents the publishing a synchronous tunnel that can be used to send messages to a consumer and wait for a response.
type Publisher interface {
	// Publish publishes a message to the tunnel, synchronously waiting for a response from the other side.
	// If no consumer is active, ErrNoConsumerActive is returned.
	Publish(ctx context.Context, message Message) (*Response, error)
}

type RedisPublishClient interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
}

type RedisPublisher struct {
	redisClient RedisPublishClient
	channelName string
}

// NewRedisPublisher creates a new RedisPublisher that publishes messages to the given channel.
func NewRedisPublisher(redisClient RedisPublishClient, channelName string) *RedisPublisher {
	return &RedisPublisher{redisClient: redisClient, channelName: channelName}
}

func (p *RedisPublisher) Publish(ctx context.Context, message *Message) (*Response, error) {
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
		return nil, ErrNoActiveConsumer
	}

	// Wait for a response from the consumer.
	responseKey := GetResponseKeyForMessage(message.ID)
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

	var response *Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}
