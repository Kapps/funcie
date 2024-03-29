// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	funcie "github.com/Kapps/funcie/pkg/funcie"
	mock "github.com/stretchr/testify/mock"
)

// ApplicationClient is an autogenerated mock type for the ApplicationClient type
type ApplicationClient struct {
	mock.Mock
}

type ApplicationClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ApplicationClient) EXPECT() *ApplicationClient_Expecter {
	return &ApplicationClient_Expecter{mock: &_m.Mock}
}

// ProcessRequest provides a mock function with given fields: ctx, application, request
func (_m *ApplicationClient) ProcessRequest(ctx context.Context, application funcie.Application, request *funcie.Message) (*funcie.Response, error) {
	ret := _m.Called(ctx, application, request)

	var r0 *funcie.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, funcie.Application, *funcie.Message) (*funcie.Response, error)); ok {
		return rf(ctx, application, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, funcie.Application, *funcie.Message) *funcie.Response); ok {
		r0 = rf(ctx, application, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, funcie.Application, *funcie.Message) error); ok {
		r1 = rf(ctx, application, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ApplicationClient_ProcessRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessRequest'
type ApplicationClient_ProcessRequest_Call struct {
	*mock.Call
}

// ProcessRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - application funcie.Application
//   - request *funcie.Message
func (_e *ApplicationClient_Expecter) ProcessRequest(ctx interface{}, application interface{}, request interface{}) *ApplicationClient_ProcessRequest_Call {
	return &ApplicationClient_ProcessRequest_Call{Call: _e.mock.On("ProcessRequest", ctx, application, request)}
}

func (_c *ApplicationClient_ProcessRequest_Call) Run(run func(ctx context.Context, application funcie.Application, request *funcie.Message)) *ApplicationClient_ProcessRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(funcie.Application), args[2].(*funcie.Message))
	})
	return _c
}

func (_c *ApplicationClient_ProcessRequest_Call) Return(_a0 *funcie.Response, _a1 error) *ApplicationClient_ProcessRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ApplicationClient_ProcessRequest_Call) RunAndReturn(run func(context.Context, funcie.Application, *funcie.Message) (*funcie.Response, error)) *ApplicationClient_ProcessRequest_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewApplicationClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewApplicationClient creates a new instance of ApplicationClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplicationClient(t mockConstructorTestingTNewApplicationClient) *ApplicationClient {
	mock := &ApplicationClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
