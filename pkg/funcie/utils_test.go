package funcie_test

import (
	. "github.com/Kapps/funcie/pkg/funcie"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsRunningWithLambda(t *testing.T) {
	t.Run("should return true if the environment variable is set", func(t *testing.T) {
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test")

		runningWithLambda := IsRunningWithLambda()

		require.True(t, runningWithLambda)
	})

	t.Run("should return false if the environment variable is not set", func(t *testing.T) {
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")

		runningWithLambda := IsRunningWithLambda()

		require.False(t, runningWithLambda)
	})
}
