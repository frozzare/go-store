package redis

import (
	"testing"

	"github.com/frozzare/go-assert"
)

/*
func TestCustomOptions(t *testing.T) {
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
	s, _ := Open()

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", "Fredrik")

	v, _ = s.Get("name")
	assert.Equal(t, "Fredrik", v.(string))

	s.Delete("name")
}

func TestGetSetMap(t *testing.T) {
	s, _ := Open()

	v, _ := s.Get("map")
	assert.Nil(t, v)

	s.Set("map", map[string]interface{}{"hello": "world"})

	v, _ = s.Get("map")
	assert.Equal(t, "world", v.(map[string]interface{})["hello"].(string))

	s.Delete("map")
}

func TestCount(t *testing.T) {
	s, _ := Open()

	c, _ := s.Count()
	assert.Equal(t, 0, c)

	s.Set("name", []byte("Fredrik"))
	c, _ = s.Count()
	assert.Equal(t, 1, c)

	s.Delete("name")
}

func TestExists(t *testing.T) {
	s, _ := Open()

	e, _ := s.Exists("name")
	assert.False(t, e)

	s.Set("name", []byte("Fredrik"))
	e, _ = s.Exists("name")
	assert.True(t, e)

	s.Delete("name")
}

func TestDeleteSimple(t *testing.T) {
	s, _ := Open()

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
	s, _ := Open()

	k, _ := s.Keys()
	assert.Equal(t, 0, len(k))

	s.Set("name", "Fredrik")

	k, _ = s.Keys()

	assert.Equal(t, 1, len(k))
	assert.Equal(t, "name", k[0])

	s.Delete("name")
}

func TestFlush(t *testing.T) {
	s, _ := Open()

	s.Set("name", "Fredrik")

	c, _ := s.Count()
	assert.Equal(t, 1, c)

	assert.Nil(t, s.Flush())

	c, _ = s.Count()
	assert.Equal(t, 0, c)
}

type Person struct {
	Name string
}

func TestGetSetSimpleStruct(t *testing.T) {
	s, _ := Open()

	v, _ := s.Get("name")
	assert.Nil(t, v)

	s.Set("name", &Person{Name: "Fredrik"})

	var p *Person
	s.Get("name", &p)
	assert.Equal(t, "Fredrik", p.Name)

	s.Delete("name")
}
