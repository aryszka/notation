package notation

import (
	"os"
	"strings"
	"testing"
)

func withEnv(t *testing.T, e ...string) (revert func()) {
	var r []func()
	revert = func() {
		for i := range r {
			r[i]()
		}
	}

	revertOne := func(key, value string, existed bool) func() {
		return func() {
			if existed {
				if err := os.Setenv(key, value); err != nil {
					t.Fatal(err)
				}

				return
			}

			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}

	for i := range e {
		var key, value string
		p := strings.Split(e[i], "=")
		key = p[0]
		if len(p) > 1 {
			value = p[1]
		}

		prev, ok := os.LookupEnv(key)
		if err := os.Setenv(key, value); err != nil {
			revert()
			t.Fatal(err)
		}

		r = append(r, revertOne(key, prev, ok))
	}

	return
}
