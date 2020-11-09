package notation

import (
	"bytes"
	"errors"
	"testing"
)

type failingWriter int

var errTest = errors.New("test")

func failAfter(n int) *failingWriter {
	w := failingWriter(n)
	return &w
}

func (w *failingWriter) Write(p []byte) (int, error) {
	*w = failingWriter(int(*w) - len(p))
	if *w >= 0 {
		return len(p), nil
	}

	return len(p) + int(*w), errTest
}

func TestFailingWriter(t *testing.T) {
	t.Run("single object", func(t *testing.T) {
		o := struct{ fooBarBaz int }{42}
		w := failAfter(9)
		n, err := Fprint(w, o)
		if n == 9 && err == errTest {
			return
		}

		if n != 9 {
			t.Fatalf("failed to writ the expected bytes; expected: 9, written: %d", n)
		}

		if err != errTest {
			t.Fatalf("failed to receive the right error; expected: %v, received: %v", errTest, err)
		}
	})

	t.Run("multiple objects, fail first", func(t *testing.T) {
		o := struct{ fooBarBaz int }{42}
		w := failAfter(9)
		n, err := Fprint(w, o, o)
		if n == 9 && err == errTest {
			return
		}

		if n != 9 {
			t.Fatalf("failed to writ the expected bytes; expected: 9, written: %d", n)
		}

		if err != errTest {
			t.Fatalf("failed to receive the right error; expected: %v, received: %v", errTest, err)
		}
	})

	t.Run("multiple objects, fail second", func(t *testing.T) {
		o := struct{ fooBarBaz int }{42}
		w := failAfter(18)
		n, err := Fprint(w, o, o)
		if n == 9 && err == errTest {
			return
		}

		if n != 18 {
			t.Fatalf("failed to writ the expected bytes; expected: 9, written: %d", n)
		}

		if err != errTest {
			t.Fatalf("failed to receive the right error; expected: %v, received: %v", errTest, err)
		}
	})
}

func TestFprint(t *testing.T) {
	t.Run("Fprint", func(t *testing.T) {
		const expect = `{fooBarBaz: 42}`
		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprint(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})

	t.Run("Fprintw", func(t *testing.T) {
		const expect = `{
	fooBarBaz: 42,
}`

		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprintw(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})

	t.Run("Fprintt", func(t *testing.T) {
		const expect = `struct{fooBarBaz int}{fooBarBaz: 42}`
		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprintt(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})

	t.Run("Fprintwt", func(t *testing.T) {
		const expect = `struct{
	fooBarBaz int
}{
	fooBarBaz: 42,
}`

		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprintwt(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})

	t.Run("Fprintv", func(t *testing.T) {
		const expect = `struct{fooBarBaz int}{fooBarBaz: int(42)}`
		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprintv(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})

	t.Run("Fprintv", func(t *testing.T) {
		const expect = `struct{
	fooBarBaz int
}{
	fooBarBaz: int(42),
}`
		var b bytes.Buffer
		o := struct{ fooBarBaz int }{42}
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		n, err := Fprintwv(&b, o)
		if err != nil {
			t.Fatal(err)
		}

		if n != len(expect) {
			t.Fatalf("invalid write length; expected: %d, got: %d", len(expect), n)
		}

		if b.String() != expect {
			t.Fatalf("invalid output; expected: %s, got: %s", expect, b.String())
		}
	})
}
