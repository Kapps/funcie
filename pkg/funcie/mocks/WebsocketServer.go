// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	publisher "github.com/Kapps/funcie/pkg/funcie/transports/ws/publisher"
	mock "github.com/stretchr/testify/mock"

	websocket "nhooyr.io/websocket"
)

// WebsocketServer is an autogenerated mock type for the WebsocketServer type
type WebsocketServer struct {
	mock.Mock
}

type WebsocketServer_Expecter struct {
	mock *mock.Mock
}

func (_m *WebsocketServer) EXPECT() *WebsocketServer_Expecter {
	return &WebsocketServer_Expecter{mock: &_m.Mock}
}

// Accept provides a mock function with given fields: w, r, opts
func (_m *WebsocketServer) Accept(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions) (publisher.Websocket, error) {
	ret := _m.Called(w, r, opts)

	var r0 publisher.Websocket
	var r1 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request, *websocket.AcceptOptions) (publisher.Websocket, error)); ok {
		return rf(w, r, opts)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request, *websocket.AcceptOptions) publisher.Websocket); ok {
		r0 = rf(w, r, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(publisher.Websocket)
		}
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request, *websocket.AcceptOptions) error); ok {
		r1 = rf(w, r, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WebsocketServer_Accept_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Accept'
type WebsocketServer_Accept_Call struct {
	*mock.Call
}

// Accept is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - r *http.Request
//   - opts *websocket.AcceptOptions
func (_e *WebsocketServer_Expecter) Accept(w interface{}, r interface{}, opts interface{}) *WebsocketServer_Accept_Call {
	return &WebsocketServer_Accept_Call{Call: _e.mock.On("Accept", w, r, opts)}
}

func (_c *WebsocketServer_Accept_Call) Run(run func(w http.ResponseWriter, r *http.Request, opts *websocket.AcceptOptions)) *WebsocketServer_Accept_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request), args[2].(*websocket.AcceptOptions))
	})
	return _c
}

func (_c *WebsocketServer_Accept_Call) Return(_a0 publisher.Websocket, _a1 error) *WebsocketServer_Accept_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *WebsocketServer_Accept_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request, *websocket.AcceptOptions) (publisher.Websocket, error)) *WebsocketServer_Accept_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewWebsocketServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebsocketServer creates a new instance of WebsocketServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebsocketServer(t mockConstructorTestingTNewWebsocketServer) *WebsocketServer {
	mock := &WebsocketServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}