// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	redis "github.com/redis/go-redis/v9"
	mock "github.com/stretchr/testify/mock"
)

// PubSub is an autogenerated mock type for the PubSub type
type PubSub struct {
	mock.Mock
}

type PubSub_Expecter struct {
	mock *mock.Mock
}

func (_m *PubSub) EXPECT() *PubSub_Expecter {
	return &PubSub_Expecter{mock: &_m.Mock}
}

// Channel provides a mock function with given fields: opts
func (_m *PubSub) Channel(opts ...redis.ChannelOption) <-chan *redis.Message {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 <-chan *redis.Message
	if rf, ok := ret.Get(0).(func(...redis.ChannelOption) <-chan *redis.Message); ok {
		r0 = rf(opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *redis.Message)
		}
	}

	return r0
}

// PubSub_Channel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Channel'
type PubSub_Channel_Call struct {
	*mock.Call
}

// Channel is a helper method to define mock.On call
//   - opts ...redis.ChannelOption
func (_e *PubSub_Expecter) Channel(opts ...interface{}) *PubSub_Channel_Call {
	return &PubSub_Channel_Call{Call: _e.mock.On("Channel",
		append([]interface{}{}, opts...)...)}
}

func (_c *PubSub_Channel_Call) Run(run func(opts ...redis.ChannelOption)) *PubSub_Channel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]redis.ChannelOption, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(redis.ChannelOption)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *PubSub_Channel_Call) Return(_a0 <-chan *redis.Message) *PubSub_Channel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PubSub_Channel_Call) RunAndReturn(run func(...redis.ChannelOption) <-chan *redis.Message) *PubSub_Channel_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *PubSub) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PubSub_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type PubSub_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *PubSub_Expecter) Close() *PubSub_Close_Call {
	return &PubSub_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *PubSub_Close_Call) Run(run func()) *PubSub_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PubSub_Close_Call) Return(_a0 error) *PubSub_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PubSub_Close_Call) RunAndReturn(run func() error) *PubSub_Close_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewPubSub interface {
	mock.TestingT
	Cleanup(func())
}

// NewPubSub creates a new instance of PubSub. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPubSub(t mockConstructorTestingTNewPubSub) *PubSub {
	mock := &PubSub{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
