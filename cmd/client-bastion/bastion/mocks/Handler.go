// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	funcie "github.com/Kapps/funcie/pkg/funcie"
	messages "github.com/Kapps/funcie/pkg/funcie/messages"

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

// Deregister provides a mock function with given fields: ctx, message
func (_m *Handler) Deregister(ctx context.Context, message funcie.MessageBase[messages.DeregistrationRequestPayload]) (*funcie.ResponseBase[messages.DeregistrationResponsePayload], error) {
	ret := _m.Called(ctx, message)

	var r0 *funcie.ResponseBase[messages.DeregistrationResponsePayload]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.DeregistrationRequestPayload]) (*funcie.ResponseBase[messages.DeregistrationResponsePayload], error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.DeregistrationRequestPayload]) *funcie.ResponseBase[messages.DeregistrationResponsePayload]); ok {
		r0 = rf(ctx, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.ResponseBase[messages.DeregistrationResponsePayload])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, funcie.MessageBase[messages.DeregistrationRequestPayload]) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Handler_Deregister_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deregister'
type Handler_Deregister_Call struct {
	*mock.Call
}

// Deregister is a helper method to define mock.On call
//   - ctx context.Context
//   - message funcie.MessageBase[messages.DeregistrationRequestPayload]
func (_e *Handler_Expecter) Deregister(ctx interface{}, message interface{}) *Handler_Deregister_Call {
	return &Handler_Deregister_Call{Call: _e.mock.On("Deregister", ctx, message)}
}

func (_c *Handler_Deregister_Call) Run(run func(ctx context.Context, message funcie.MessageBase[messages.DeregistrationRequestPayload])) *Handler_Deregister_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(funcie.MessageBase[messages.DeregistrationRequestPayload]))
	})
	return _c
}

func (_c *Handler_Deregister_Call) Return(_a0 *funcie.ResponseBase[messages.DeregistrationResponsePayload], _a1 error) *Handler_Deregister_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Handler_Deregister_Call) RunAndReturn(run func(context.Context, funcie.MessageBase[messages.DeregistrationRequestPayload]) (*funcie.ResponseBase[messages.DeregistrationResponsePayload], error)) *Handler_Deregister_Call {
	_c.Call.Return(run)
	return _c
}

// ForwardRequest provides a mock function with given fields: ctx, message
func (_m *Handler) ForwardRequest(ctx context.Context, message funcie.MessageBase[messages.ForwardRequestPayload]) (*funcie.ResponseBase[messages.ForwardRequestResponsePayload], error) {
	ret := _m.Called(ctx, message)

	var r0 *funcie.ResponseBase[messages.ForwardRequestResponsePayload]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.ForwardRequestPayload]) (*funcie.ResponseBase[messages.ForwardRequestResponsePayload], error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.ForwardRequestPayload]) *funcie.ResponseBase[messages.ForwardRequestResponsePayload]); ok {
		r0 = rf(ctx, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.ResponseBase[messages.ForwardRequestResponsePayload])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, funcie.MessageBase[messages.ForwardRequestPayload]) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Handler_ForwardRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ForwardRequest'
type Handler_ForwardRequest_Call struct {
	*mock.Call
}

// ForwardRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - message funcie.MessageBase[messages.ForwardRequestPayload]
func (_e *Handler_Expecter) ForwardRequest(ctx interface{}, message interface{}) *Handler_ForwardRequest_Call {
	return &Handler_ForwardRequest_Call{Call: _e.mock.On("ForwardRequest", ctx, message)}
}

func (_c *Handler_ForwardRequest_Call) Run(run func(ctx context.Context, message funcie.MessageBase[messages.ForwardRequestPayload])) *Handler_ForwardRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(funcie.MessageBase[messages.ForwardRequestPayload]))
	})
	return _c
}

func (_c *Handler_ForwardRequest_Call) Return(_a0 *funcie.ResponseBase[messages.ForwardRequestResponsePayload], _a1 error) *Handler_ForwardRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Handler_ForwardRequest_Call) RunAndReturn(run func(context.Context, funcie.MessageBase[messages.ForwardRequestPayload]) (*funcie.ResponseBase[messages.ForwardRequestResponsePayload], error)) *Handler_ForwardRequest_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, message
func (_m *Handler) Register(ctx context.Context, message funcie.MessageBase[messages.RegistrationRequestPayload]) (*funcie.ResponseBase[messages.RegistrationResponsePayload], error) {
	ret := _m.Called(ctx, message)

	var r0 *funcie.ResponseBase[messages.RegistrationResponsePayload]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.RegistrationRequestPayload]) (*funcie.ResponseBase[messages.RegistrationResponsePayload], error)); ok {
		return rf(ctx, message)
	}
	if rf, ok := ret.Get(0).(func(context.Context, funcie.MessageBase[messages.RegistrationRequestPayload]) *funcie.ResponseBase[messages.RegistrationResponsePayload]); ok {
		r0 = rf(ctx, message)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.ResponseBase[messages.RegistrationResponsePayload])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, funcie.MessageBase[messages.RegistrationRequestPayload]) error); ok {
		r1 = rf(ctx, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Handler_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type Handler_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - message funcie.MessageBase[messages.RegistrationRequestPayload]
func (_e *Handler_Expecter) Register(ctx interface{}, message interface{}) *Handler_Register_Call {
	return &Handler_Register_Call{Call: _e.mock.On("Register", ctx, message)}
}

func (_c *Handler_Register_Call) Run(run func(ctx context.Context, message funcie.MessageBase[messages.RegistrationRequestPayload])) *Handler_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(funcie.MessageBase[messages.RegistrationRequestPayload]))
	})
	return _c
}

func (_c *Handler_Register_Call) Return(_a0 *funcie.ResponseBase[messages.RegistrationResponsePayload], _a1 error) *Handler_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Handler_Register_Call) RunAndReturn(run func(context.Context, funcie.MessageBase[messages.RegistrationRequestPayload]) (*funcie.ResponseBase[messages.RegistrationResponsePayload], error)) *Handler_Register_Call {
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
