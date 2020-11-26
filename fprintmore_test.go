package notation

import (
	"bytes"
	"errors"
	"io"
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
	defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
	o := struct{ fooBarBaz int }{42}
	for _, test := range []struct {
		name   string
		fn     func(io.Writer, ...interface{}) (int, error)
		expect string
	}{{
		name:   "Fprint",
		fn:     Fprint,
		expect: `{fooBarBaz: 42}`,
	}, {
		name: "Fprintw",
		fn:   Fprintw,
		expect: `{
	fooBarBaz: 42,
}`,
	}, {
		name:   "Fprintt",
		fn:     Fprintt,
		expect: `struct{fooBarBaz int}{fooBarBaz: 42}`,
	}, {
		name: "Fprintwt",
		fn:   Fprintwt,
		expect: `struct{
	fooBarBaz int
}{
	fooBarBaz: 42,
}`,
	}, {
		name:   "Fprintv",
		fn:     Fprintv,
		expect: `struct{fooBarBaz int}{fooBarBaz: int(42)}`,
	}, {
		name: "Fprintwv",
		fn:   Fprintwv,
		expect: `struct{
	fooBarBaz int
}{
	fooBarBaz: int(42),
}`,
	}} {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			n, err := test.fn(&b, o)
			if err != nil {
				t.Fatal(err)
			}

			if n != len(test.expect) {
				t.Fatalf("invalid write length; expected %d, got: %d", len(test.expect), n)
			}

			if b.String() != test.expect {
				t.Fatalf("invalid output; expected: %s, got: %s", test.expect, b.String())
			}
		})
	}
}
