package funcie_test

import (
	"encoding/json"
	"errors"
	. "github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResponseUnmarshal(t *testing.T) {
	t.Parallel()

	response := NewResponse("id", []byte("data"), nil)
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var resp Response
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, response, &resp)
}

func TestResponseUnmarshal_WithError(t *testing.T) {
	t.Parallel()

	response := NewResponse("id", nil, errors.New("error"))
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var resp Response
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, response, &resp)
}
