package tunnel_test

import (
	"errors"
	. "funcie/tunnel"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewProxyError(t *testing.T) {
	proxyError := NewProxyError("test")
	require.Equal(t, "test", proxyError.Message)
}

func TestProxyError_Error(t *testing.T) {
	proxyError := NewProxyError("test")
	require.Equal(t, "test", proxyError.Error())
}

func TestProxyError_MarshalJSON(t *testing.T) {
	proxyError := NewProxyError("test")
	json, err := proxyError.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "{\"message\":\"test\"}", string(json))
}

func TestProxyError_UnmarshalJSON(t *testing.T) {
	proxyError := &ProxyError{}
	err := proxyError.UnmarshalJSON([]byte("{\"message\":\"test\"}"))
	require.NoError(t, err)
	require.Equal(t, "test", proxyError.Message)
}

func TestProxyError_UnmarshalJSON_InvalidJson(t *testing.T) {
	proxyError := &ProxyError{}
	err := proxyError.UnmarshalJSON([]byte("invalid json"))
	require.Error(t, err)
}

func TestNewProxyErrorFromError(t *testing.T) {
	proxyError := NewProxyErrorFromError(nil)
	require.Nil(t, proxyError)
}

func TestNewProxyErrorFromError_WithError(t *testing.T) {
	proxyError := NewProxyErrorFromError(errors.New("test"))
	require.Equal(t, "test", proxyError.Message)
}
