package tunnel

import (
	"encoding/json"
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

// RoughCompareMatcher returns a matcher that compares the given object to the
// expected object using cmp.Diff. If the objects are not equal, the matcher
// will return false; otherwise, it will return true.
func RoughCompareMatcher[T comparable](
	expected T,
) func(T) bool {
	return func(actual T) bool {
		return RoughCompare(expected, actual)
	}
}

// RoughCompareMatcherJson returns a matcher that compares the given object to the
// input object, when the expected object is a JSON string.
// The object is unmarshalled from JSON, and then compared using RoughCompare.
// If the unmarshalling fails, or the objects are not equal, the matcher
// will return false; otherwise, it will return true.
func RoughCompareMatcherJson[T comparable](
	expected T,
) func(string) bool {
	return func(actualJson string) bool {
		var actual T
		err := json.Unmarshal([]byte(actualJson), &actual)
		if err != nil {
			return false
		}

		return RoughCompare(expected, actual)
	}
}
