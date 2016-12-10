package redis

import (
	"testing"

	"github.com/frozzare/go-assert"
)

/*
func TestCustomClient(t *testing.T) {
	s := Open(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
}*/

func TestGetSetSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
}

func TestGetSetMap(t *testing.T) {
	s := Open()

	v, _ := s.Get("map")
	assert.Nil(t, v)

	s.Set("map", map[string]interface{}{"hello": "world"})

	v, _ = s.Get("map")
	assert.Equal(t, "world", v.(map[string]interface{})["hello"].(string))

	s.Delete("map")
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

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
	v, _ = s.Get("name")
	assert.Nil(t, v)
}
