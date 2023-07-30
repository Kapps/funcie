// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// BastionReceiver is an autogenerated mock type for the BastionReceiver type
type BastionReceiver struct {
	mock.Mock
}

type BastionReceiver_Expecter struct {
	mock *mock.Mock
}

func (_m *BastionReceiver) EXPECT() *BastionReceiver_Expecter {
	return &BastionReceiver_Expecter{mock: &_m.Mock}
}

// Start provides a mock function with given fields:
func (_m *BastionReceiver) Start() {
	_m.Called()
}

// BastionReceiver_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type BastionReceiver_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
func (_e *BastionReceiver_Expecter) Start() *BastionReceiver_Start_Call {
	return &BastionReceiver_Start_Call{Call: _e.mock.On("Start")}
}

func (_c *BastionReceiver_Start_Call) Run(run func()) *BastionReceiver_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BastionReceiver_Start_Call) Return() *BastionReceiver_Start_Call {
	_c.Call.Return()
	return _c
}

func (_c *BastionReceiver_Start_Call) RunAndReturn(run func()) *BastionReceiver_Start_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *BastionReceiver) Stop() {
	_m.Called()
}

// BastionReceiver_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type BastionReceiver_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *BastionReceiver_Expecter) Stop() *BastionReceiver_Stop_Call {
	return &BastionReceiver_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *BastionReceiver_Stop_Call) Run(run func()) *BastionReceiver_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *BastionReceiver_Stop_Call) Return() *BastionReceiver_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *BastionReceiver_Stop_Call) RunAndReturn(run func()) *BastionReceiver_Stop_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewBastionReceiver interface {
	mock.TestingT
	Cleanup(func())
}

// NewBastionReceiver creates a new instance of BastionReceiver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBastionReceiver(t mockConstructorTestingTNewBastionReceiver) *BastionReceiver {
	mock := &BastionReceiver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
