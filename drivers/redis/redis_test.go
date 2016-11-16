package redis

import (
	"testing"

	redis "gopkg.in/redis.v3"

	assert "github.com/frozzare/go-assert"
)

func TestCustomClient(t *testing.T) {
	s := Open(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}))

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", []byte("Fredrik"))

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", string(v))

	s.Delete("name")
}

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
