// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	funcie "github.com/Kapps/funcie/pkg/funcie"
	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

type Handler_Expecter struct {
	mock *mock.Mock
}

func (_m *Handler) EXPECT() *Handler_Expecter {
	return &Handler_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, message
func (_m *Handler) Execute(ctx context.Context, message *funcie.Message) (*funcie.Response, error) {
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

// Handler_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Handler_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - message *funcie.Message
func (_e *Handler_Expecter) Execute(ctx interface{}, message interface{}) *Handler_Execute_Call {
	return &Handler_Execute_Call{Call: _e.mock.On("Execute", ctx, message)}
}

func (_c *Handler_Execute_Call) Run(run func(ctx context.Context, message *funcie.Message)) *Handler_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*funcie.Message))
	})
	return _c
}

func (_c *Handler_Execute_Call) Return(_a0 *funcie.Response, _a1 error) *Handler_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Handler_Execute_Call) RunAndReturn(run func(context.Context, *funcie.Message) (*funcie.Response, error)) *Handler_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHandler(t mockConstructorTestingTNewHandler) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
