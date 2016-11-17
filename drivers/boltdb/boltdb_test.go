package boltdb

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestCustomClient(t *testing.T) {
	s := Open("/tmp/custom.db")

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", []byte("Fredrik"))

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name")
}

func TestGetSetSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name1")
	assert.Nil(t, v)

	s.Set("name1", []byte("Fredrik"))

	v, _ = s.Get("name1")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name1")
}

func TestCount(t *testing.T) {
	s := Open()

	assert.Equal(t, 0, s.Count())

	s.Set("name2", []byte("Fredrik"))
	assert.Equal(t, 1, s.Count())

	s.Delete("name2")
}

func TestExists(t *testing.T) {
	s := Open()

	assert.False(t, s.Exists("name3"))

	s.Set("name3", []byte("Fredrik"))
	assert.True(t, s.Exists("name3"))

	s.Delete("name3")
}

func TestDeleteSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name4")
	assert.Nil(t, v)

	s.Set("name4", []byte("Fredrik"))

	v, _ = s.Get("name4")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name4")
	v, _ = s.Get("name4")
	assert.Nil(t, v)
}
