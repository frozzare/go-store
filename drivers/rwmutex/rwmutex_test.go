package rwmutex

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestGetSetSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", []byte("Fredrik"))

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name")
}

func TestCount(t *testing.T) {
	s := Open()

	assert.Equal(t, 0, s.Count())

	s.Set("name", []byte("Fredrik"))
	assert.Equal(t, 1, s.Count())

	s.Delete("name")
}

func TestExists(t *testing.T) {
	s := Open()

	assert.False(t, s.Exists("name"))

	s.Set("name", []byte("Fredrik"))
	assert.True(t, s.Exists("name"))

	s.Delete("name")
}

func TestDeleteSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", []byte("Fredrik"))

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name")
	v, _ = s.Get("name")
	assert.Nil(t, v)
}

func TestInstance(t *testing.T) {
	assert.Equal(t, 0, Open().Count())
	assert.Equal(t, 0, Open("cache").Count())

	Open("cache").Set("name", []byte("Fredrik"))

	assert.Equal(t, 0, Open().Count())
	assert.Equal(t, 1, Open("cache").Count())
}
