package funcie_test

import (
	"github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewApplication(t *testing.T) {
	t.Parallel()

	endpoint := funcie.NewEndpoint("http", "host", 1234)
	application := funcie.NewApplication("name", endpoint)
	require.Equal(t, "name", application.Name)
	require.Equal(t, endpoint, application.Endpoint)
}

func TestApplication_String(t *testing.T) {
	t.Parallel()

	endpoint := funcie.NewEndpoint("http", "host", 1234)
	application := funcie.NewApplication("name", endpoint)
	require.Equal(t, "name (http://host:1234)", application.String())
}
