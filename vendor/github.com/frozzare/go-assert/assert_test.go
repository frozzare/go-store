package assert

import "testing"

type Hello struct {
	Name string
}

func TestEqual(t *testing.T) {
	Equal(t, "foo", "foo")

	p := Hello{"Fredrik"}
	p2 := Hello{"Kalle"}

	NotEqual(t, p, p2)
}

func TestTrue(t *testing.T) {
	True(t, true)
}

func TestFalse(t *testing.T) {
	False(t, false)
}

func TestNil(t *testing.T) {
	Nil(t, nil)
	var i interface{}
	Nil(t, i)
}

func TestNotNil(t *testing.T) {
	NotNil(t, true)

	p := Hello{"Fredrik"}
	NotNil(t, p)
}
