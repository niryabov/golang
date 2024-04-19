//go:build !solution

package testequal

import (
	"fmt"
	"regexp"
	"slices"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but contin	ues execution.
//
// Returns true iff arguments are equal.

func Report(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if len(msgAndArgs) != 0 {
		if len(msgAndArgs) > 1 {
			t.Errorf("not equal\nexpected: %v\nactual: %v\nmessage:%v", fmt.Sprint(expected), fmt.Sprint(actual), fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...))
		} else {
			t.Errorf("not equal\nexpected: %v\nactual: %v\nmessage:%v", fmt.Sprint(expected), fmt.Sprint(actual), msgAndArgs[0].(string))

		}
	} else {
		t.Errorf("not equal\nexpected: %v\nactual: %v", fmt.Sprint(expected), fmt.Sprint(actual))
	}
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	var flag bool
	if fmt.Sprintf("%T", expected) != fmt.Sprintf("%T", actual) {
		t.Errorf("Unmatched arguments' types")
		return false
	}

	switch actual.(type) {
	case map[string]string:
		if (expected.(map[string]string) == nil) != (actual.(map[string]string) == nil) {
			Report(t, expected, actual, msgAndArgs...)
			return false
		}
		flag = fmt.Sprint(expected) == fmt.Sprint(actual)

	case []int:
		if (expected.([]int) == nil) != (actual.([]int) == nil) {
			return false
		}
		flag = slices.Equal(expected.([]int), actual.([]int))
	case []byte:
		if (expected.([]byte) == nil) != (actual.([]byte) == nil) {
			return false
		}
		flag = slices.Equal(expected.([]byte), actual.([]byte))
	default:
		matched, _ := regexp.MatchString(`int*|uint*|string`, fmt.Sprintf("%T", actual))
		if !matched {
			t.Errorf("Unmatched arguments' types")
			return false
		}
		flag = (expected == actual)
	}

	if !flag {
		Report(t, expected, actual, msgAndArgs...)
	}
	return flag
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	var flag bool
	if fmt.Sprintf("%T", expected) != fmt.Sprintf("%T", actual) {
		return true
	}

	switch actual.(type) {
	case map[string]string:
		if (expected.(map[string]string) == nil) != (actual.(map[string]string) == nil) {
			return true
		}
		flag = fmt.Sprint(expected) == fmt.Sprint(actual)

	case []int:
		if (expected.([]int) == nil) != (actual.([]int) == nil) {
			return true
		}
		flag = slices.Equal(expected.([]int), actual.([]int))
	case []byte:
		if (expected.([]byte) == nil) != (actual.([]byte) == nil) {
			return true
		}
		flag = slices.Equal(expected.([]byte), actual.([]byte))

	default:
		matched, _ := regexp.MatchString(`int*|uint*|string`, fmt.Sprintf("%T", actual))
		if !matched {
			return true
		}

		flag = (expected == actual)
	}

	if flag {
		Report(t, expected, actual, msgAndArgs...)
	}
	return !flag
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !AssertEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !AssertNotEqual(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}
