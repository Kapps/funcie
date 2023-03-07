// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	redis "github.com/redis/go-redis/v9"

	time "time"
)

// RedisPublishClient is an autogenerated mock type for the RedisPublishClient type
type RedisPublishClient struct {
	mock.Mock
}

type RedisPublishClient_Expecter struct {
	mock *mock.Mock
}

func (_m *RedisPublishClient) EXPECT() *RedisPublishClient_Expecter {
	return &RedisPublishClient_Expecter{mock: &_m.Mock}
}

// BRPop provides a mock function with given fields: ctx, timeout, keys
func (_m *RedisPublishClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	_va := make([]interface{}, len(keys))
	for _i := range keys {
		_va[_i] = keys[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, timeout)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *redis.StringSliceCmd
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, ...string) *redis.StringSliceCmd); ok {
		r0 = rf(ctx, timeout, keys...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.StringSliceCmd)
		}
	}

	return r0
}

// RedisPublishClient_BRPop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BRPop'
type RedisPublishClient_BRPop_Call struct {
	*mock.Call
}

// BRPop is a helper method to define mock.On call
//   - ctx context.Context
//   - timeout time.Duration
//   - keys ...string
func (_e *RedisPublishClient_Expecter) BRPop(ctx interface{}, timeout interface{}, keys ...interface{}) *RedisPublishClient_BRPop_Call {
	return &RedisPublishClient_BRPop_Call{Call: _e.mock.On("BRPop",
		append([]interface{}{ctx, timeout}, keys...)...)}
}

func (_c *RedisPublishClient_BRPop_Call) Run(run func(ctx context.Context, timeout time.Duration, keys ...string)) *RedisPublishClient_BRPop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(time.Duration), variadicArgs...)
	})
	return _c
}

func (_c *RedisPublishClient_BRPop_Call) Return(_a0 *redis.StringSliceCmd) *RedisPublishClient_BRPop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RedisPublishClient_BRPop_Call) RunAndReturn(run func(context.Context, time.Duration, ...string) *redis.StringSliceCmd) *RedisPublishClient_BRPop_Call {
	_c.Call.Return(run)
	return _c
}

// Publish provides a mock function with given fields: ctx, channel, message
func (_m *RedisPublishClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	ret := _m.Called(ctx, channel, message)

	var r0 *redis.IntCmd
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) *redis.IntCmd); ok {
		r0 = rf(ctx, channel, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.IntCmd)
		}
	}

	return r0
}

// RedisPublishClient_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type RedisPublishClient_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - ctx context.Context
//   - channel string
//   - message interface{}
func (_e *RedisPublishClient_Expecter) Publish(ctx interface{}, channel interface{}, message interface{}) *RedisPublishClient_Publish_Call {
	return &RedisPublishClient_Publish_Call{Call: _e.mock.On("Publish", ctx, channel, message)}
}

func (_c *RedisPublishClient_Publish_Call) Run(run func(ctx context.Context, channel string, message interface{})) *RedisPublishClient_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *RedisPublishClient_Publish_Call) Return(_a0 *redis.IntCmd) *RedisPublishClient_Publish_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RedisPublishClient_Publish_Call) RunAndReturn(run func(context.Context, string, interface{}) *redis.IntCmd) *RedisPublishClient_Publish_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRedisPublishClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewRedisPublishClient creates a new instance of RedisPublishClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRedisPublishClient(t mockConstructorTestingTNewRedisPublishClient) *RedisPublishClient {
	mock := &RedisPublishClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}