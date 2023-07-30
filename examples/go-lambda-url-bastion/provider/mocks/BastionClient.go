// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"
	json "encoding/json"

	funcie "github.com/Kapps/funcie/pkg/funcie"

	mock "github.com/stretchr/testify/mock"
)

// BastionClient is an autogenerated mock type for the BastionClient type
type BastionClient struct {
	mock.Mock
}

type BastionClient_Expecter struct {
	mock *mock.Mock
}

func (_m *BastionClient) EXPECT() *BastionClient_Expecter {
	return &BastionClient_Expecter{mock: &_m.Mock}
}

// SendRequest provides a mock function with given fields: ctx, request
func (_m *BastionClient) SendRequest(ctx context.Context, request *funcie.MessageBase[json.RawMessage]) (*funcie.ResponseBase[json.RawMessage], error) {
	ret := _m.Called(ctx, request)

	var r0 *funcie.ResponseBase[json.RawMessage]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *funcie.MessageBase[json.RawMessage]) (*funcie.ResponseBase[json.RawMessage], error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *funcie.MessageBase[json.RawMessage]) *funcie.ResponseBase[json.RawMessage]); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*funcie.ResponseBase[json.RawMessage])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *funcie.MessageBase[json.RawMessage]) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BastionClient_SendRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendRequest'
type BastionClient_SendRequest_Call struct {
	*mock.Call
}

// SendRequest is a helper method to define mock.On call
//   - ctx context.Context
//   - request *funcie.MessageBase[json.RawMessage]
func (_e *BastionClient_Expecter) SendRequest(ctx interface{}, request interface{}) *BastionClient_SendRequest_Call {
	return &BastionClient_SendRequest_Call{Call: _e.mock.On("SendRequest", ctx, request)}
}

func (_c *BastionClient_SendRequest_Call) Run(run func(ctx context.Context, request *funcie.MessageBase[json.RawMessage])) *BastionClient_SendRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*funcie.MessageBase[json.RawMessage]))
	})
	return _c
}

func (_c *BastionClient_SendRequest_Call) Return(_a0 *funcie.ResponseBase[json.RawMessage], _a1 error) *BastionClient_SendRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BastionClient_SendRequest_Call) RunAndReturn(run func(context.Context, *funcie.MessageBase[json.RawMessage]) (*funcie.ResponseBase[json.RawMessage], error)) *BastionClient_SendRequest_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewBastionClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewBastionClient creates a new instance of BastionClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBastionClient(t mockConstructorTestingTNewBastionClient) *BastionClient {
	mock := &BastionClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}