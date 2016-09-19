package store

import (
	"testing"

	"github.com/frozzare/go-assert"
)

func TestGetSetSimple(t *testing.T) {
	s := New()

	assert.Nil(t, s.Get("name"))

	s.Set("name", "Fredrik")
	assert.Equal(t, "Fredrik", s.Get("name"))
}

func TestGetSetMap(t *testing.T) {
	s := New()

	assert.Nil(t, s.Get("map"))

	s.Set("map", map[string]interface{}{
		"name": "Fredrik",
	})
	assert.Equal(t, "Fredrik", s.Get("map").(map[string]interface{})["name"])
}

func TestCount(t *testing.T) {
	s := New()

	assert.Equal(t, 0, s.Count())

	s.Set("name", "Fredrik")
	assert.Equal(t, 1, s.Count())
}

func TestExists(t *testing.T) {
	s := New()

	assert.False(t, s.Exists("name"))

	s.Set("name", "Fredrik")
	assert.True(t, s.Exists("name"))
}

func TestDeleteSimple(t *testing.T) {
	s := New()

	assert.Nil(t, s.Get("name"))

	s.Set("name", "Fredrik")
	assert.Equal(t, "Fredrik", s.Get("name"))

	s.Delete("name")
	assert.Nil(t, s.Get("name"))
}

func TestInstance(t *testing.T) {
	assert.Equal(t, 0, Instance().Count())
	assert.Equal(t, 0, Instance("cache").Count())

	Instance("cache").Set("name", "Fredrik")

	assert.Equal(t, 0, Instance().Count())
	assert.Equal(t, 1, Instance("cache").Count())
}
