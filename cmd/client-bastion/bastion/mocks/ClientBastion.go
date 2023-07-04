// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ClientBastion is an autogenerated mock type for the ClientBastion type
type ClientBastion struct {
	mock.Mock
}

type ClientBastion_Expecter struct {
	mock *mock.Mock
}

func (_m *ClientBastion) EXPECT() *ClientBastion_Expecter {
	return &ClientBastion_Expecter{mock: &_m.Mock}
}

// Listen provides a mock function with given fields:
func (_m *ClientBastion) Listen() {
	_m.Called()
}

// ClientBastion_Listen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Listen'
type ClientBastion_Listen_Call struct {
	*mock.Call
}

// Listen is a helper method to define mock.On call
func (_e *ClientBastion_Expecter) Listen() *ClientBastion_Listen_Call {
	return &ClientBastion_Listen_Call{Call: _e.mock.On("Listen")}
}

func (_c *ClientBastion_Listen_Call) Run(run func()) *ClientBastion_Listen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ClientBastion_Listen_Call) Return() *ClientBastion_Listen_Call {
	_c.Call.Return()
	return _c
}

func (_c *ClientBastion_Listen_Call) RunAndReturn(run func()) *ClientBastion_Listen_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewClientBastion interface {
	mock.TestingT
	Cleanup(func())
}

// NewClientBastion creates a new instance of ClientBastion. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClientBastion(t mockConstructorTestingTNewClientBastion) *ClientBastion {
	mock := &ClientBastion{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}