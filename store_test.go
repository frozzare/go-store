package store

import (
	"testing"

	assert "github.com/frozzare/go-assert"
	"github.com/frozzare/go-store/rwmutex"
)

func TestOpenNoDriver(t *testing.T) {
	driver, err := Open("test")
	assert.Nil(t, driver)
	assert.NotNil(t, err)
}

func TestOpenDriver(t *testing.T) {
	driver, err := Open("rwmutex")
	assert.NotNil(t, driver)
	assert.Nil(t, err)

	// "rwmutex" is default value.
	driver, err = Open()
	assert.NotNil(t, driver)
	assert.Nil(t, err)
}

func TestRegisterNilDriver(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Register("test", nil)
}

func TestRegisterSameDriver(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	Register("rwmutex", &rwmutex.Driver{})
}
