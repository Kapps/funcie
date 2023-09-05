package receiver

import (
	"context"
	"fmt"
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/redis/go-redis/v9"
)

// RedisClient is an interface for the redis.Client type containing only the methods we use.
type RedisClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
}

type redisApplicationRegistry struct {
	redisClient RedisClient
}

var appKeyBase = "funcie:apps"

// NewRedisApplicationRegistry creates a new Redis-backed application registry.
func NewRedisApplicationRegistry(redisClient RedisClient) funcie.ApplicationRegistry {
	return &redisApplicationRegistry{redisClient: redisClient}
}

func (r *redisApplicationRegistry) Register(ctx context.Context, application *funcie.Application) error {
	key := getKeyForApplication(application.Name)
	res := r.redisClient.HSet(ctx, key, "endpoint", application.Endpoint.String())
	if err := res.Err(); err != nil {
		return fmt.Errorf("register application: %w", err)
	}

	return nil
}

func (r *redisApplicationRegistry) Unregister(ctx context.Context, applicationName string) error {
	key := getKeyForApplication(applicationName)

	res := r.redisClient.Del(ctx, key)
	if err := res.Err(); err != nil {
		return fmt.Errorf("unregister application: %w", err)
	}

	return nil
}

func (r *redisApplicationRegistry) GetApplication(ctx context.Context, applicationName string) (*funcie.Application, error) {
	key := getKeyForApplication(applicationName)

	res := r.redisClient.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("getting application with key %v: %w", key, err)
	}

	vals := res.Val()
	if len(vals) == 0 {
		return nil, funcie.ErrApplicationNotFound
	}

	endpoint, err := funcie.NewEndpointFromAddress(vals["endpoint"])
	if err != nil {
		return nil, fmt.Errorf("parsing endpoint %v: %w", vals["endpoint"], err)
	}

	return &funcie.Application{
		Name:     applicationName,
		Endpoint: endpoint,
	}, nil
}

func getKeyForApplication(applicationName string) string {
	return fmt.Sprintf("%s:%s", appKeyBase, applicationName)
}
