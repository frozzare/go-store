package rethinkdb

import (
	"testing"

	assert "github.com/frozzare/go-assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestCustomOptions(t *testing.T) {
	s := Open(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "custom",
	}, "custom")

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
}

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

func TestKeys(t *testing.T) {
	s := Open()

	k, _ := s.Keys()
	assert.Equal(t, 0, len(k))

	s.Set("name", "Fredrik")

	k, _ = s.Keys()

	assert.Equal(t, 1, len(k))
	assert.Equal(t, "name", k[0])

	s.Delete("name")
}
