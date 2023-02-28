// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	tunnel "funcie/tunnel"

	mock "github.com/stretchr/testify/mock"
)

// Publisher is an autogenerated mock type for the Publisher type
type Publisher struct {
	mock.Mock
}

type Publisher_Expecter struct {
	mock *mock.Mock
}

func (_m *Publisher) EXPECT() *Publisher_Expecter {
	return &Publisher_Expecter{mock: &_m.Mock}
}

// Publish provides a mock function with given fields: ctx, message
func (_m *Publisher) Publish(ctx context.Context, message tunnel.Message) (*tunnel.Response, error) {
	ret := _m.Called(ctx, message)

	var r0 *tunnel.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, tunnel.Message) (*tunnel.Response, error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, tunnel.Message) *tunnel.Response); ok {
		r0 = rf(ctx, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tunnel.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, tunnel.Message) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publisher_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type Publisher_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - ctx context.Context
//   - message tunnel.Message
func (_e *Publisher_Expecter) Publish(ctx interface{}, message interface{}) *Publisher_Publish_Call {
	return &Publisher_Publish_Call{Call: _e.mock.On("Publish", ctx, message)}
}

func (_c *Publisher_Publish_Call) Run(run func(ctx context.Context, message tunnel.Message)) *Publisher_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(tunnel.Message))
	})
	return _c
}

func (_c *Publisher_Publish_Call) Return(_a0 *tunnel.Response, _a1 error) *Publisher_Publish_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Publisher_Publish_Call) RunAndReturn(run func(context.Context, tunnel.Message) (*tunnel.Response, error)) *Publisher_Publish_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewPublisher interface {
	mock.TestingT
	Cleanup(func())
}

// NewPublisher creates a new instance of Publisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPublisher(t mockConstructorTestingTNewPublisher) *Publisher {
	mock := &Publisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
