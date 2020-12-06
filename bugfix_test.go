package notation

import "testing"

func TestMapWrapping(t *testing.T) {
	const expect = `map[*interface{}]interface{}{
	*interface{}(string("foo")): interface{}(bool(true)),
	*interface{}(string("foo")): interface{}(bool(true)),
}`

	ifpointer := func(v interface{}) *interface{} {
		return &v
	}

	m := map[*interface{}]interface{}{
		ifpointer("foo"): true,
		ifpointer("foo"): true,
	}

	s := Sprintwv(m)
	if s != expect {
		t.Fatalf("expected: %s, got: %s", expect, s)
	}
}
