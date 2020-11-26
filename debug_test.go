package notation

import (
	"fmt"
	"reflect"
	"testing"
)

	func TestDebugNode(t *testing.T) {
	const expect = `"foobarbaz"`
	o := "foobarbaz"
	n := reflectValue(none, &pending{values: make(map[uintptr]nodeRef)}, reflect.ValueOf(o))
	s := fmt.Sprint(n)
	if s != expect {
		t.Fatalf(
			"failed to get debug string of node, got: %s, expected: %s",
			s,
			expect,
		)
	}
}
