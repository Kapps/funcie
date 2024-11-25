// Code generated by mockery v2.33.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	websocket "nhooyr.io/websocket"
)

// Connection is an autogenerated mock type for the Connection type
type Connection struct {
	mock.Mock
}

type Connection_Expecter struct {
	mock *mock.Mock
}

func (_m *Connection) EXPECT() *Connection_Expecter {
	return &Connection_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields: code, reason
func (_m *Connection) Close(code websocket.StatusCode, reason string) error {
	ret := _m.Called(code, reason)

	var r0 error
	if rf, ok := ret.Get(0).(func(websocket.StatusCode, string) error); ok {
		r0 = rf(code, reason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connection_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Connection_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
//   - code websocket.StatusCode
//   - reason string
func (_e *Connection_Expecter) Close(code interface{}, reason interface{}) *Connection_Close_Call {
	return &Connection_Close_Call{Call: _e.mock.On("Close", code, reason)}
}

func (_c *Connection_Close_Call) Run(run func(code websocket.StatusCode, reason string)) *Connection_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(websocket.StatusCode), args[1].(string))
	})
	return _c
}

func (_c *Connection_Close_Call) Return(_a0 error) *Connection_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Connection_Close_Call) RunAndReturn(run func(websocket.StatusCode, string) error) *Connection_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx, message
func (_m *Connection) Read(ctx context.Context, message interface{}) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connection_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type Connection_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - message interface{}
func (_e *Connection_Expecter) Read(ctx interface{}, message interface{}) *Connection_Read_Call {
	return &Connection_Read_Call{Call: _e.mock.On("Read", ctx, message)}
}

func (_c *Connection_Read_Call) Run(run func(ctx context.Context, message interface{})) *Connection_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *Connection_Read_Call) Return(_a0 error) *Connection_Read_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Connection_Read_Call) RunAndReturn(run func(context.Context, interface{}) error) *Connection_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: ctx, message
func (_m *Connection) Write(ctx context.Context, message interface{}) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connection_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type Connection_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - ctx context.Context
//   - message interface{}
func (_e *Connection_Expecter) Write(ctx interface{}, message interface{}) *Connection_Write_Call {
	return &Connection_Write_Call{Call: _e.mock.On("Write", ctx, message)}
}

func (_c *Connection_Write_Call) Run(run func(ctx context.Context, message interface{})) *Connection_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(interface{}))
	})
	return _c
}

func (_c *Connection_Write_Call) Return(_a0 error) *Connection_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Connection_Write_Call) RunAndReturn(run func(context.Context, interface{}) error) *Connection_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewConnection creates a new instance of Connection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConnection(t interface {
	mock.TestingT
	Cleanup(func())
}) *Connection {
	mock := &Connection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
