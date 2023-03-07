package funcie

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"math"
	"time"
)

// RoughCompare compares the given object to the expected object using rough comparisons.
// Some types, such as time.Time, will be checked for "roughly equal" as opposed to identity comparisons.
func RoughCompare(expected interface{}, actual interface{}) bool {
	timeComparer := cmp.Comparer(func(x, y time.Time) bool {
		return math.Abs(float64(x.Sub(y).Milliseconds())) < 500
	})
	diff := cmp.Diff(actual, expected, timeComparer)
	if diff != "" {
		return false
	}
	return true
}

// RequireRoughCompare compares the given object to the expected object using RoughCompare.
// If the objects are not equal, the test will fail.
func RequireRoughCompare(t require.TestingT, expected interface{}, actual interface{}) {
	if !RoughCompare(expected, actual) {
		t.Errorf("RoughCompare failed: %v", cmp.Diff(actual, expected))
	}
}
