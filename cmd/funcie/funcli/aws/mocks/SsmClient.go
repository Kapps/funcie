// Code generated by mockery v2.33.1. DO NOT EDIT.

package mocks

import (
	context "context"

	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	mock "github.com/stretchr/testify/mock"
)

// SsmClient is an autogenerated mock type for the SsmClient type
type SsmClient struct {
	mock.Mock
}

type SsmClient_Expecter struct {
	mock *mock.Mock
}

func (_m *SsmClient) EXPECT() *SsmClient_Expecter {
	return &SsmClient_Expecter{mock: &_m.Mock}
}

// GetParameter provides a mock function with given fields: ctx, params, optFns
func (_m *SsmClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ssm.GetParameterOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) *ssm.GetParameterOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssm.GetParameterOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SsmClient_GetParameter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetParameter'
type SsmClient_GetParameter_Call struct {
	*mock.Call
}

// GetParameter is a helper method to define mock.On call
//   - ctx context.Context
//   - params *ssm.GetParameterInput
//   - optFns ...func(*ssm.Options)
func (_e *SsmClient_Expecter) GetParameter(ctx interface{}, params interface{}, optFns ...interface{}) *SsmClient_GetParameter_Call {
	return &SsmClient_GetParameter_Call{Call: _e.mock.On("GetParameter",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *SsmClient_GetParameter_Call) Run(run func(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options))) *SsmClient_GetParameter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*ssm.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*ssm.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*ssm.GetParameterInput), variadicArgs...)
	})
	return _c
}

func (_c *SsmClient_GetParameter_Call) Return(_a0 *ssm.GetParameterOutput, _a1 error) *SsmClient_GetParameter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SsmClient_GetParameter_Call) RunAndReturn(run func(context.Context, *ssm.GetParameterInput, ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)) *SsmClient_GetParameter_Call {
	_c.Call.Return(run)
	return _c
}

// StartSession provides a mock function with given fields: ctx, params, optFns
func (_m *SsmClient) StartSession(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options)) (*ssm.StartSessionOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ssm.StartSessionOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.StartSessionInput, ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.StartSessionInput, ...func(*ssm.Options)) *ssm.StartSessionOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssm.StartSessionOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssm.StartSessionInput, ...func(*ssm.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SsmClient_StartSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartSession'
type SsmClient_StartSession_Call struct {
	*mock.Call
}

// StartSession is a helper method to define mock.On call
//   - ctx context.Context
//   - params *ssm.StartSessionInput
//   - optFns ...func(*ssm.Options)
func (_e *SsmClient_Expecter) StartSession(ctx interface{}, params interface{}, optFns ...interface{}) *SsmClient_StartSession_Call {
	return &SsmClient_StartSession_Call{Call: _e.mock.On("StartSession",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *SsmClient_StartSession_Call) Run(run func(ctx context.Context, params *ssm.StartSessionInput, optFns ...func(*ssm.Options))) *SsmClient_StartSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*ssm.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*ssm.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*ssm.StartSessionInput), variadicArgs...)
	})
	return _c
}

func (_c *SsmClient_StartSession_Call) Return(_a0 *ssm.StartSessionOutput, _a1 error) *SsmClient_StartSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SsmClient_StartSession_Call) RunAndReturn(run func(context.Context, *ssm.StartSessionInput, ...func(*ssm.Options)) (*ssm.StartSessionOutput, error)) *SsmClient_StartSession_Call {
	_c.Call.Return(run)
	return _c
}

// TerminateSession provides a mock function with given fields: ctx, params, optFns
func (_m *SsmClient) TerminateSession(ctx context.Context, params *ssm.TerminateSessionInput, optFns ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ssm.TerminateSessionOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.TerminateSessionInput, ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.TerminateSessionInput, ...func(*ssm.Options)) *ssm.TerminateSessionOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssm.TerminateSessionOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssm.TerminateSessionInput, ...func(*ssm.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SsmClient_TerminateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TerminateSession'
type SsmClient_TerminateSession_Call struct {
	*mock.Call
}

// TerminateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - params *ssm.TerminateSessionInput
//   - optFns ...func(*ssm.Options)
func (_e *SsmClient_Expecter) TerminateSession(ctx interface{}, params interface{}, optFns ...interface{}) *SsmClient_TerminateSession_Call {
	return &SsmClient_TerminateSession_Call{Call: _e.mock.On("TerminateSession",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *SsmClient_TerminateSession_Call) Run(run func(ctx context.Context, params *ssm.TerminateSessionInput, optFns ...func(*ssm.Options))) *SsmClient_TerminateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*ssm.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*ssm.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*ssm.TerminateSessionInput), variadicArgs...)
	})
	return _c
}

func (_c *SsmClient_TerminateSession_Call) Return(_a0 *ssm.TerminateSessionOutput, _a1 error) *SsmClient_TerminateSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SsmClient_TerminateSession_Call) RunAndReturn(run func(context.Context, *ssm.TerminateSessionInput, ...func(*ssm.Options)) (*ssm.TerminateSessionOutput, error)) *SsmClient_TerminateSession_Call {
	_c.Call.Return(run)
	return _c
}

// NewSsmClient creates a new instance of SsmClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSsmClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SsmClient {
	mock := &SsmClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}