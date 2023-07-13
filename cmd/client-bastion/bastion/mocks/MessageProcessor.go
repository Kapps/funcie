// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	funcie "github.com/Kapps/funcie/pkg/funcie"
	mock "github.com/stretchr/testify/mock"
)

// MessageProcessor is an autogenerated mock type for the MessageProcessor type
type MessageProcessor struct {
	mock.Mock
}

type MessageProcessor_Expecter struct {
	mock *mock.Mock
}

func (_m *MessageProcessor) EXPECT() *MessageProcessor_Expecter {
	return &MessageProcessor_Expecter{mock: &_m.Mock}
}

// ProcessMessage provides a mock function with given fields: ctx, message
func (_m *MessageProcessor) ProcessMessage(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
	ret := _m.Called(ctx, message)

	var r0 *funcie.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *funcie.Message) (*funcie.Response, error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *funcie.Message) *funcie.Response); ok {
		r0 = rf(ctx, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *funcie.Message) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageProcessor_ProcessMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessMessage'
type MessageProcessor_ProcessMessage_Call struct {
	*mock.Call
}

// ProcessMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - message *funcie.Message
func (_e *MessageProcessor_Expecter) ProcessMessage(ctx interface{}, message interface{}) *MessageProcessor_ProcessMessage_Call {
	return &MessageProcessor_ProcessMessage_Call{Call: _e.mock.On("ProcessMessage", ctx, message)}
}

func (_c *MessageProcessor_ProcessMessage_Call) Run(run func(ctx context.Context, message *funcie.Message)) *MessageProcessor_ProcessMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*funcie.Message))
	})
	return _c
}

func (_c *MessageProcessor_ProcessMessage_Call) Return(_a0 *funcie.Response, _a1 error) *MessageProcessor_ProcessMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MessageProcessor_ProcessMessage_Call) RunAndReturn(run func(context.Context, *funcie.Message) (*funcie.Response, error)) *MessageProcessor_ProcessMessage_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMessageProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageProcessor creates a new instance of MessageProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageProcessor(t mockConstructorTestingTNewMessageProcessor) *MessageProcessor {
	mock := &MessageProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}