// Code generated by mockery v2.21.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	websocket "nhooyr.io/websocket"
)

// Websocket is an autogenerated mock type for the Websocket type
type Websocket struct {
	mock.Mock
}

type Websocket_Expecter struct {
	mock *mock.Mock
}

func (_m *Websocket) EXPECT() *Websocket_Expecter {
	return &Websocket_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields: code, reason
func (_m *Websocket) Close(code websocket.StatusCode, reason string) error {
	ret := _m.Called(code, reason)

	var r0 error
	if rf, ok := ret.Get(0).(func(websocket.StatusCode, string) error); ok {
		r0 = rf(code, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Websocket_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Websocket_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - code websocket.StatusCode
//   - reason string
func (_e *Websocket_Expecter) Close(code interface{}, reason interface{}) *Websocket_Close_Call {
	return &Websocket_Close_Call{Call: _e.mock.On("Close", code, reason)}
}

func (_c *Websocket_Close_Call) Run(run func(code websocket.StatusCode, reason string)) *Websocket_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(websocket.StatusCode), args[1].(string))
	})
	return _c
}

func (_c *Websocket_Close_Call) Return(_a0 error) *Websocket_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Websocket_Close_Call) RunAndReturn(run func(websocket.StatusCode, string) error) *Websocket_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx
func (_m *Websocket) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	ret := _m.Called(ctx)

	var r0 websocket.MessageType
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) (websocket.MessageType, []byte, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) websocket.MessageType); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(websocket.MessageType)
	}

	if rf, ok := ret.Get(1).(func(context.Context) []byte); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Websocket_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type Websocket_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Websocket_Expecter) Read(ctx interface{}) *Websocket_Read_Call {
	return &Websocket_Read_Call{Call: _e.mock.On("Read", ctx)}
}

func (_c *Websocket_Read_Call) Run(run func(ctx context.Context)) *Websocket_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Websocket_Read_Call) Return(_a0 websocket.MessageType, _a1 []byte, _a2 error) *Websocket_Read_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Websocket_Read_Call) RunAndReturn(run func(context.Context) (websocket.MessageType, []byte, error)) *Websocket_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: ctx, typ, p
func (_m *Websocket) Write(ctx context.Context, typ websocket.MessageType, p []byte) error {
	ret := _m.Called(ctx, typ, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, websocket.MessageType, []byte) error); ok {
		r0 = rf(ctx, typ, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Websocket_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type Websocket_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - ctx context.Context
//   - typ websocket.MessageType
//   - p []byte
func (_e *Websocket_Expecter) Write(ctx interface{}, typ interface{}, p interface{}) *Websocket_Write_Call {
	return &Websocket_Write_Call{Call: _e.mock.On("Write", ctx, typ, p)}
}

func (_c *Websocket_Write_Call) Run(run func(ctx context.Context, typ websocket.MessageType, p []byte)) *Websocket_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(websocket.MessageType), args[2].([]byte))
	})
	return _c
}

func (_c *Websocket_Write_Call) Return(_a0 error) *Websocket_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Websocket_Write_Call) RunAndReturn(run func(context.Context, websocket.MessageType, []byte) error) *Websocket_Write_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewWebsocket interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebsocket creates a new instance of Websocket. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebsocket(t mockConstructorTestingTNewWebsocket) *Websocket {
	mock := &Websocket{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}