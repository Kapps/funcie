package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/Kapps/funcie/pkg/funcie/transports/utils"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

// PubSub is the interface that wraps the redis PubSub methods used by the consumer.
type PubSub interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
	Receive(ctx context.Context) (interface{}, error)
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
	redisClient     ConsumeClient
	pubsub          PubSub
	router          utils.ClientHandlerRouter
	baseChannelName string
}

// NewConsumer creates a new RedisConsumer that consumes messages from channels starting with the given base name.
// This implementation takes in a *redis.Client instead of a RedisConsumeClient so that it can be used with a real Redis client.
func NewConsumer(redisClient *redis.Client, baseChannelName string, router utils.ClientHandlerRouter) funcie.Consumer {
	wrappedRedis := &redisConsumeClient{
		Client: redisClient,
	}
	return &Consumer{
		redisClient:     wrappedRedis,
		baseChannelName: baseChannelName,
		router:          router,
	}
}

// NewConsumerWithClient creates a new RedisConsumer that consumes messages from channels starting with the given base name.
func NewConsumerWithClient(redisClient ConsumeClient, baseChannelName string, router utils.ClientHandlerRouter) funcie.Consumer {
	return &Consumer{
		redisClient:     redisClient,
		baseChannelName: baseChannelName,
		router:          router,
	}
}

func (c *Consumer) Connect(ctx context.Context) error {
	// To connect, we can just subscribe to the base channel name.
	ps := c.redisClient.Subscribe(ctx, c.baseChannelName)
	received, err := ps.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive from pubsub: %w", err)
	}

	sub := received.(*redis.Subscription)
	slog.Info("subscribed to base channel", "channel", sub.Channel, "count", sub.Count)

	c.pubsub = ps
	return nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	ps := c.pubsub
	defer funcie.CloseOrLog(fmt.Sprintf("pubsub from base channel %v", c.baseChannelName), ps)

	slog.InfoContext(ctx, "starting to consume messages", "baseChannelName", c.baseChannelName)

	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled", "err", ctx.Err())
			return ctx.Err()
		case msg, ok := <-ps.Channel():
			if !ok {
				slog.Debug("pubsub channel closed")
				return funcie.ErrPubSubChannelClosed
			}

			slog.DebugContext(ctx, "received message", "channel", msg.Channel)
			go func(msg *redis.Message) {
				err := c.processMessage(ctx, msg)
				if err != nil {
					// If we get an error processing the message, we still want to continue our loop.
					// So we just log the error and keep going.
					slog.ErrorContext(ctx, "error processing message", err)
				}
			}(msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg *redis.Message) error {
	slog.DebugContext(ctx, "received message", "channel", msg.Channel, "payload", msg.Payload)

	message, err := parseMessage(msg.Payload)
	if err != nil {
		return fmt.Errorf("error parsing message: %w", err)
	}

	response, err := c.router.Handle(ctx, message)
	if err == utils.ErrNoHandlerFound {
		unsubErr := c.Unsubscribe(ctx, message.Application)
		if unsubErr != nil {
			// An error unsubscribing isn't the end of the world. We can still continue and still want to return the original error.
			slog.ErrorContext(ctx, "error unsubscribing from channel", err, "channel", msg.Channel)
		}
		return fmt.Errorf("no handler found for app %v in message %v: %w", message.Application, message.ID, err)
	}
	if err != nil {
		return fmt.Errorf("error handling message: %w", err)
	}

	responseKey := GetResponseKeyForMessage(c.baseChannelName, message.ID)
	responseData, err := formatResponse(response)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	cmd := c.redisClient.RPush(ctx, responseKey, responseData)
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("error pushing response to queue: %w", err)
	}

	return nil
}

func (c *Consumer) Subscribe(ctx context.Context, applicationId string, handler funcie.Handler) error {
	channelName := GetChannelNameForApplication(c.baseChannelName, applicationId)
	slog.Info("subscribing to channel", "channel", channelName)

	if err := c.pubsub.Subscribe(ctx, channelName); err != nil {
		return fmt.Errorf("subscribing to channel: %w", err)
	}

	if err := c.router.AddClientHandler(applicationId, handler); err != nil {
		return fmt.Errorf("adding client handler: %w", err)
	}

	return nil
}

func (c *Consumer) Unsubscribe(ctx context.Context, applicationId string) error {
	channelName := GetChannelNameForApplication(c.baseChannelName, applicationId)
	slog.Info("unsubscribing from channel", "channel", channelName)

	if err := c.router.RemoveClientHandler(applicationId); err != nil {
		return fmt.Errorf("removing client handler: %w", err)
	}

	if err := c.pubsub.Unsubscribe(ctx, channelName); err != nil {
		return fmt.Errorf("unsubscribing from channel: %w", err)
	}

	return nil
}

func parseMessage(message string) (*funcie.Message, error) {
	var msg funcie.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return nil, fmt.Errorf("unmarshalling message: %w", err)
	}
	return &msg, nil
}

func formatResponse(response *funcie.Response) (string, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("marshalling response: %w", err)
	}
	return string(data), nil
}
