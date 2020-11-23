package notation

import (
	"bytes"
	"testing"
)

func TestPrint(t *testing.T) {
	for _, test := range []struct {
		p func(...interface{}) (int, error)
		o interface{}
		e string
	}{{
		p: Print,
		o: struct{ foo int }{42},
		e: `{foo: 42}`,
	}, {
		p: Printw,
		o: struct{ foo int }{42},
		e: `{
	foo: 42,
}`,
	}, {
		p: Printt,
		o: struct{ foo int }{42},
		e: `struct{foo int}{foo: 42}`,
	}, {
		p: Printwt,
		o: struct{ foo int }{42},
		e: `struct{
	foo int
}{
	foo: 42,
}`,
	}, {
		p: Printv,
		o: struct{ foo int }{42},
		e: `struct{foo int}{foo: int(42)}`,
	}, {
		p: Printwv,
		o: struct{ foo int }{42},
		e: `struct{
	foo int
}{
	foo: int(42),
}`,
	}} {
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		t.Run("", func(t *testing.T) {
			var b bytes.Buffer
			orig := stderr
			stderr = &b
			defer func() { stderr = orig }()
			n, err := test.p(test.o)
			if err != nil {
				t.Fatal(err)
			}

			if n != len(test.e) {
				t.Fatalf("expected length: %d, got: %d", len(test.e), n)
			}

			s := b.String()
			if s != test.e {
				t.Fatalf("expected: %s, got: %s", test.e, s)
			}
		})
	}
}

func TestPrintln(t *testing.T) {
	for _, test := range []struct {
		p func(...interface{}) (int, error)
		f bool
		o interface{}
		e string
	}{{
		p: Println,
		o: struct{ foo int }{42},
		e: "{foo: 42}\n",
	}, {
		p: Printlnw,
		o: struct{ foo int }{42},
		e: `{
	foo: 42,
}
`,
	}, {
		p: Printlnt,
		o: struct{ foo int }{42},
		e: "struct{foo int}{foo: 42}\n",
	}, {
		p: Printlnwt,
		o: struct{ foo int }{42},
		e: `struct{
	foo int
}{
	foo: 42,
}
`,
	}, {
		p: Printlnv,
		o: struct{ foo int }{42},
		e: "struct{foo int}{foo: int(42)}\n",
	}, {
		p: Printlnwv,
		o: struct{ foo int }{42},
		e: `struct{
	foo int
}{
	foo: int(42),
}
`,
	}, {
		p: Println,
		o: struct{ foo int }{42},
		f: true,
	}} {
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		t.Run("", func(t *testing.T) {
			var b bytes.Buffer
			orig := stderr
			stderr = &b
			if test.f {
				var w failingWriter
				stderr = &w
			}

			defer func() { stderr = orig }()
			n, err := test.p(test.o)
			if test.f && err == nil {
				t.Fatal("failed to fail")
			}

			if test.f {
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if n != len(test.e) {
				t.Fatalf("expected length: %d, got: %d", len(test.e), n)
			}

			s := b.String()
			if s != test.e {
				t.Fatalf("expected: %s, got: %s", test.e, s)
			}
		})
	}
}
