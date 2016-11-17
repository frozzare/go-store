package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// equal will test if the expected and actual value is the same match.
func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
		return true
	}

	return false
}

// Fail will print a failing message to the terminal.
func Fail(t *testing.T, expected, actual interface{}) {
	_, file, line, _ := runtime.Caller(2)

	t.Errorf("\033[31m✖\033[39m %s:%d: \033[31m%v == %v\033[39m\n",
		filepath.Base(file),
		line,
		expected,
		actual)
}

// Equal will test if the expected and actual value is the same match.
func Equal(t *testing.T, expected, actual interface{}) {
	if !equal(expected, actual) {
		Fail(t, expected, actual)
	}
}

// NotEqual will test if the expected and actual value is not a match.
func NotEqual(t *testing.T, expected, actual interface{}) {
	if equal(expected, actual) {
		Fail(t, expected, actual)
	}
}

// True will test the actual value and see if it's true or not.
func True(t *testing.T, actual bool) {
	if actual != true {
		Fail(t, true, actual)
	}
}

// False will test the actual value and see if it's false or not.
func False(t *testing.T, actual bool) {
	if actual != false {
		Fail(t, false, actual)
	}
}

// Nil will test the actual value and see if it's nil or not.
func Nil(t *testing.T, actual interface{}, args ...interface{}) {
	success := true

	if actual == nil {
		success = false
	}

	value := reflect.ValueOf(actual)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		success = false
	}

	if success {
		Fail(t, nil, actual)
	}
}

// NotNil test check if the actual value is not nil.
func NotNil(t *testing.T, actual interface{}) {
	success := true

	if actual == nil {
		success = false
	} else {
		value := reflect.ValueOf(actual)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			success = false
		}
	}

	if !success {
		Fail(t, nil, actual)
	}
}
