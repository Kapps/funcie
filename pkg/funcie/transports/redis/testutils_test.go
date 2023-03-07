package redis_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
	"time"
)

var defaultTimeout = 10 * time.Second

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

// ExpectSendToChannel expects to send a value to the given channel within 1 second.
func ExpectSendToChannel[T any](t *testing.T, channel chan<- T, value T) {
	select {
	case channel <- value:
	case <-time.After(defaultTimeout):
		t.Errorf("Expected to send to channel, but timed out")
	}
}

// ExpectReceiveFromChannel expects to receive a value from the given channel within 1 second.
func ExpectReceiveFromChannel[T any](t *testing.T, channel <-chan T) T {
	select {
	case actual := <-channel:
		return actual
	case <-time.After(defaultTimeout):
		t.Errorf("Expected to receive from channel, but timed out")
		return *new(T)
	}
}
