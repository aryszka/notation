package notation

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestSprint(t *testing.T) {
	type (
		myBool          bool
		myInt           int
		myFloat         float64
		myComplex       complex64
		myArray         [3]int
		myChannel       chan int
		myFunction      func(int, int) int
		myMap           map[int]int
		myPointer       *int
		myList          []int
		myString        string
		myStruct        struct{ field interface{} }
		myUnsafePointer unsafe.Pointer
	)

	for _, test := range []struct {
		title  string
		value  interface{}
		expect string
	}{
		{"nil", nil, "nil"},
		{"false", false, "false"},
		{"true", true, "true"},
		{"custom false", myBool(false), "false"},
		{"custom true", myBool(true), "true"},
		{"int", 42, "42"},
		{"negative int", -42, "-42"},
		{"custom int", myInt(42), "42"},
		{"uint", uint(42), "42"},
		{"byte", byte(42), "42"},
		{"round float", float64(42), "42"},
		{"custom round float", myFloat(42), "42"},
		{"float with fraction", 1.8, "1.8"},
		{"custom float with fraction", myFloat(1.8), "1.8"},
		{"complex", 2 + 3i, "2+3i"},
		{"custom complex", myComplex(2 + 3i), "2+3i"},
		{"imaginary", 3i, "0+3i"},
		{"custom imaginary", myComplex(3i), "0+3i"},
		{"array", [...]int{1, 2, 3}, "[3]{1, 2, 3}"},
		{"custom array", myArray{1, 2, 3}, "[3]{1, 2, 3}"},
		{"channel", make(chan int), "chan"},
		{"custom channel", make(myChannel), "chan"},
		{"nil channel", struct{ c chan int }{}, "{c: nil}"},
		{"nil custom channel", struct{ c myChannel }{}, "{c: nil}"},
		{"receive channel", struct{ c <-chan int }{make(chan int)}, "{c: chan}"},
		{"send channel", struct{ c chan<- int }{make(chan int)}, "{c: chan}"},
		{"function", func() {}, "func()"},
		{"function with args", func(int, int) int { return 0 }, "func()"},
		{"custom function with args", myFunction(func(int, int) int { return 0 }), "func()"},
		{"function with multiple return args", func(int, int) (int, int) { return 0, 0 }, "func()"},
		{"nil function", struct{ f func(int) int }{}, "{f: nil}"},
		{"interface type", struct{ i interface{ foo() } }{}, "{i: nil}"},
		{"map", map[int]int{24: 42}, "map{24: 42}"},
		{"custom map", myMap{24: 42}, "map{24: 42}"},
		{"nil map", struct{ m map[int]int }{}, "{m: nil}"},
		{"nil custom map", struct{ m myMap }{}, "{m: nil}"},
		{"pointer", &struct{}{}, "{}"},
		{"custom pointer", &myStruct{}, "{field: nil}"},
		{"nil pointer", struct{ p *int }{}, "{p: nil}"},
		{"nil custom pointer", struct{ p myPointer }{}, "{p: nil}"},
		{"list", []int{1, 2, 3}, "[]{1, 2, 3}"},
		{"custom list", myList{1, 2, 3}, "[]{1, 2, 3}"},
		{"nil list", struct{ l []int }{}, "{l: nil}"},
		{"nil custom list", struct{ l myList }{}, "{l: nil}"},
		{"string", "\\\"\b\f\n\r\t\vfoo", "\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\""},
		{"custom string", myString("\\\"\b\f\n\r\t\vfoo"), "\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\""},
		{"structure", struct{ foo int }{42}, "{foo: 42}"},
		{"custom structure", myStruct{42}, "{field: 42}"},
		{"unsafe pointer", unsafe.Pointer(&struct{}{}), "pointer"},
		{"custom unsafe pointer", myUnsafePointer(&struct{}{}), "pointer"},
		{"unsafe pointer type", struct{ p unsafe.Pointer }{}, "{p: nil}"},
	} {
		t.Run(test.title, func(t *testing.T) {
			s := Sprint(test.value)
			if s != test.expect {
				t.Fatalf("expected: %s, got: %s", test.expect, s)
			}
		})
	}
}

func TestSprintt(t *testing.T) {
	type (
		myBool          bool
		myInt           int
		myFloat         float64
		myComplex       complex64
		myArray         [3]int
		myChannel       chan int
		myFunction      func(int, int) int
		myMap           map[int]int
		myPointer       *int
		myList          []int
		myString        string
		myStruct        struct{ field interface{} }
		myUnsafePointer unsafe.Pointer
	)

	for _, test := range []struct {
		title  string
		value  interface{}
		expect string
	}{
		{"nil", nil, "nil"},
		{"false", false, "false"},
		{"true", true, "true"},
		{"custom false", myBool(false), "myBool(false)"},
		{"custom true", myBool(true), "myBool(true)"},
		{"int", 42, "42"},
		{"negative int", -42, "-42"},
		{"custom int", myInt(42), "myInt(42)"},
		{"uint", uint(42), "uint(42)"},
		{"byte", byte(42), "uint8(42)"},
		{"round float", float64(42), "float64(42)"},
		{"custom round float", myFloat(42), "myFloat(42)"},
		{"float with fraction", 1.8, "float64(1.8)"},
		{"custom float with fraction", myFloat(1.8), "myFloat(1.8)"},
		{"complex", 2 + 3i, "complex128(2+3i)"},
		{"custom complex", myComplex(2 + 3i), "myComplex(2+3i)"},
		{"imaginary", 3i, "complex128(0+3i)"},
		{"custom imaginary", myComplex(3i), "myComplex(0+3i)"},
		{"array", [...]int{1, 2, 3}, "[3]int{1, 2, 3}"},
		{"custom array", myArray{1, 2, 3}, "myArray{1, 2, 3}"},
		{"channel", make(chan int), "chan int"},
		{"custom channel", make(myChannel), "myChannel"},
		{"nil channel", struct{ c chan int }{}, "struct{c chan int}{c: nil}"},
		{"nil custom channel", struct{ c myChannel }{}, "struct{c myChannel}{c: nil}"},
		{"receive channel", struct{ c <-chan int }{make(chan int)}, "struct{c <-chan int}{c: chan}"},
		{"send channel", struct{ c chan<- int }{make(chan int)}, "struct{c chan<- int}{c: chan}"},
		{"function", func() {}, "func()"},
		{"function with args", func(int, int) int { return 0 }, "func(int, int) int"},
		{"custom function with args", myFunction(func(int, int) int { return 0 }), "myFunction"},
		{"function with multiple return args", func(int, int) (int, int) { return 0, 0 }, "func(int, int) (int, int)"},
		{"nil function", struct{ f func(int) int }{}, "struct{f func(int) int}{f: nil}"},
		{"interface type", struct{ i interface{ foo() } }{}, "struct{i interface{foo()}}{i: nil}"},
		{"map", map[int]int{24: 42}, "map[int]int{24: 42}"},
		{"custom map", myMap{24: 42}, "myMap{24: 42}"},
		{"nil map", struct{ m map[int]int }{}, "struct{m map[int]int}{m: nil}"},
		{"nil custom map", struct{ m myMap }{}, "struct{m myMap}{m: nil}"},
		{"pointer", &struct{}{}, "*struct{}{}"},
		{"custom pointer", &myStruct{}, "*myStruct{field: nil}"},
		{"nil pointer", struct{ p *int }{}, "struct{p *int}{p: nil}"},
		{"nil custom pointer", struct{ p myPointer }{}, "struct{p myPointer}{p: nil}"},
		{"list", []int{1, 2, 3}, "[]int{1, 2, 3}"},
		{"custom list", myList{1, 2, 3}, "myList{1, 2, 3}"},
		{"nil list", struct{ l []int }{}, "struct{l []int}{l: nil}"},
		{"nil custom list", struct{ l myList }{}, "struct{l myList}{l: nil}"},
		{"string", "\\\"\b\f\n\r\t\vfoo", "\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\""},
		{"custom string", myString("\\\"\b\f\n\r\t\vfoo"), "myString(\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\")"},
		{"structure", struct{ foo int }{42}, "struct{foo int}{foo: 42}"},
		{"custom structure", myStruct{42}, "myStruct{field: 42}"},
		{"unsafe pointer", unsafe.Pointer(&struct{}{}), "pointer"},
		{"custom unsafe pointer", myUnsafePointer(&struct{}{}), "pointer"},
		{"unsafe pointer type", struct{ p unsafe.Pointer }{}, "struct{p Pointer}{p: nil}"},
	} {
		t.Run(test.title, func(t *testing.T) {
			s := Sprintt(test.value)
			if s != test.expect {
				t.Fatalf("expected: %s, got: %s", test.expect, s)
			}
		})
	}
}

func TestSprintv(t *testing.T) {
	type (
		myBool          bool
		myInt           int
		myFloat         float64
		myComplex       complex64
		myArray         [3]int
		myChannel       chan int
		myFunction      func(int, int) int
		myMap           map[int]int
		myPointer       *int
		myList          []int
		myString        string
		myStruct        struct{ field interface{} }
		myUnsafePointer unsafe.Pointer
	)

	for _, test := range []struct {
		title  string
		value  interface{}
		expect string
	}{
		{"nil", nil, "nil"},
		{"false", false, "bool(false)"},
		{"true", true, "bool(true)"},
		{"custom false", myBool(false), "myBool(false)"},
		{"custom true", myBool(true), "myBool(true)"},
		{"int", 42, "int(42)"},
		{"negative int", -42, "int(-42)"},
		{"custom int", myInt(42), "myInt(42)"},
		{"uint", uint(42), "uint(42)"},
		{"byte", byte(42), "uint8(42)"},
		{"round float", float64(42), "float64(42)"},
		{"custom round float", myFloat(42), "myFloat(42)"},
		{"float with fraction", 1.8, "float64(1.8)"},
		{"custom float with fraction", myFloat(1.8), "myFloat(1.8)"},
		{"complex", 2 + 3i, "complex128(2+3i)"},
		{"custom complex", myComplex(2 + 3i), "myComplex(2+3i)"},
		{"imaginary", 3i, "complex128(0+3i)"},
		{"custom imaginary", myComplex(3i), "myComplex(0+3i)"},
		{"array", [...]int{1, 2, 3}, "[3]int{int(1), int(2), int(3)}"},
		{"custom array", myArray{1, 2, 3}, "myArray{int(1), int(2), int(3)}"},
		{"channel", make(chan int), "chan int"},
		{"custom channel", make(myChannel), "myChannel"},
		{"nil channel", struct{ c chan int }{}, "struct{c chan int}{c: chan int(nil)}"},
		{"nil custom channel", struct{ c myChannel }{}, "struct{c myChannel}{c: myChannel(nil)}"},
		{"receive channel", struct{ c <-chan int }{make(chan int)}, "struct{c <-chan int}{c: <-chan int}"},
		{"send channel", struct{ c chan<- int }{make(chan int)}, "struct{c chan<- int}{c: chan<- int}"},
		{"function", func() {}, "func()"},
		{"function with args", func(int, int) int { return 0 }, "func(int, int) int"},
		{"custom function with args", myFunction(func(int, int) int { return 0 }), "myFunction"},
		{"function with multiple return args", func(int, int) (int, int) { return 0, 0 }, "func(int, int) (int, int)"},
		{"nil function", struct{ f func(int) int }{}, "struct{f func(int) int}{f: func(int) int(nil)}"},
		{"interface type", struct{ i interface{ foo() } }{}, "struct{i interface{foo()}}{i: interface{foo()}(nil)}"},
		{"map", map[int]int{24: 42}, "map[int]int{int(24): int(42)}"},
		{"custom map", myMap{24: 42}, "myMap{int(24): int(42)}"},
		{"nil map", struct{ m map[int]int }{}, "struct{m map[int]int}{m: map[int]int(nil)}"},
		{"nil custom map", struct{ m myMap }{}, "struct{m myMap}{m: myMap(nil)}"},
		{"pointer", &struct{}{}, "*struct{}{}"},
		{"custom pointer", &myStruct{}, "*myStruct{field: interface{}(nil)}"},
		{"nil pointer", struct{ p *int }{}, "struct{p *int}{p: *int(nil)}"},
		{"nil custom pointer", struct{ p myPointer }{}, "struct{p myPointer}{p: myPointer(nil)}"},
		{"list", []int{1, 2, 3}, "[]int{int(1), int(2), int(3)}"},
		{"custom list", myList{1, 2, 3}, "myList{int(1), int(2), int(3)}"},
		{"nil list", struct{ l []int }{}, "struct{l []int}{l: []int(nil)}"},
		{"nil custom list", struct{ l myList }{}, "struct{l myList}{l: myList(nil)}"},
		{"string", "\\\"\b\f\n\r\t\vfoo", "string(\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\")"},
		{"custom string", myString("\\\"\b\f\n\r\t\vfoo"), "myString(\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\")"},
		{"structure", struct{ foo int }{42}, "struct{foo int}{foo: int(42)}"},
		{"custom structure", myStruct{42}, "myStruct{field: interface{}(int(42))}"},
		{"custom structure, nil field", myStruct{}, "myStruct{field: interface{}(nil)}"},
		{"unsafe pointer", unsafe.Pointer(&struct{}{}), "pointer"},
		{"custom unsafe pointer", myUnsafePointer(&struct{}{}), "pointer"},
		{"unsafe pointer type", struct{ p unsafe.Pointer }{}, "struct{p Pointer}{p: Pointer(nil)}"},
	} {
		t.Run(test.title, func(t *testing.T) {
			s := Sprintv(test.value)
			if s != test.expect {
				t.Fatalf("expected: %s, got: %s", test.expect, s)
			}
		})
	}
}

func TestSprintInvalid(t *testing.T) {
	s := sprint(none, reflect.Value{})
	if s != "<invalid>" {
		t.Fatalf("expected: <invalid>, got: %s", s)
	}
}
