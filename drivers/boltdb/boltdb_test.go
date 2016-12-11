package boltdb

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestCustomOptions(t *testing.T) {
	s := Open("/tmp/custom-boltdb.db")

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
}

func TestGetSetSimple(t *testing.T) {
	s := Open()

	v, _ := s.Get("name1")
	assert.Nil(t, v)

	s.Set("name1", "Fredrik")

	v, _ = s.Get("name1")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name1")
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

	s.Set("name4", "Fredrik")

	v, _ = s.Get("name4")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name4")
	v, _ = s.Get("name4")
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
