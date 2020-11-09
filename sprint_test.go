package notation

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

type test struct {
	title  string
	value  interface{}
	expect string
}

type tests []test

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

func (test test) run(t *testing.T, sprint func(...interface{}) string) {
	t.Run(test.title, func(t *testing.T) {
		s := sprint(test.value)
		if s != test.expect {
			t.Fatalf("expected: %s, got: %s", test.expect, s)
		}
	})
}

func (tests tests) run(t *testing.T, sprint func(...interface{}) string) {
	logEnv := func(name string, dflt int) {
		t.Logf("%s=%d", name, config(name, dflt))
	}

	logEnv("TABWIDTH", 8)
	logEnv("LINEWIDTH", 80-8)
	logEnv("LINEWIDTH1", 80*3/2-8)

	for _, ti := range tests {
		ti.run(t, sprint)
	}
}

func (t tests) expect(expect map[string]string) tests {
	var set tests
	for _, test := range t {
		if expect, doSet := expect[test.title]; doSet {
			test.expect = expect
		}

		set = append(set, test)
	}

	return set
}

func expectedMaxWidth(t []test) int {
	var w int
	for _, test := range t {
		if len(test.expect) > w {
			w = len(test.expect)
		}
	}

	return w
}

func defaultSet() tests {
	return []test{
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
		{"long item", []string{"foobarbazqux"}, "[]{\"foobarbazqux\"}"},
		{"long subitem", []struct{ foo string }{{foo: "foobarbazqux"}}, "[]{{foo: \"foobarbazqux\"}}"},
		{"string", "\\\"\b\f\n\r\t\vfoo", "\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\""},
		{"custom string", myString("\\\"\b\f\n\r\t\vfoo"), "\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\""},
		{"structure", struct{ foo int }{42}, "{foo: 42}"},
		{"custom structure", myStruct{42}, "{field: 42}"},
		{"unsafe pointer", unsafe.Pointer(&struct{}{}), "pointer"},
		{"custom unsafe pointer", myUnsafePointer(&struct{}{}), "pointer"},
		{"unsafe pointer type", struct{ p unsafe.Pointer }{}, "{p: nil}"},
	}
}

func (t tests) expectTypes() tests {
	return t.expect(map[string]string{
		"custom false":                       "myBool(false)",
		"custom true":                        "myBool(true)",
		"custom int":                         "myInt(42)",
		"uint":                               "uint(42)",
		"byte":                               "byte(42)",
		"round float":                        "float64(42)",
		"custom round float":                 "myFloat(42)",
		"float with fraction":                "float64(1.8)",
		"custom float with fraction":         "myFloat(1.8)",
		"complex":                            "complex128(2+3i)",
		"custom complex":                     "myComplex(2+3i)",
		"imaginary":                          "complex128(0+3i)",
		"custom imaginary":                   "myComplex(0+3i)",
		"array":                              "[3]int{1, 2, 3}",
		"custom array":                       "myArray{1, 2, 3}",
		"channel":                            "chan int",
		"custom channel":                     "myChannel",
		"nil channel":                        "struct{c chan int}{c: nil}",
		"nil custom channel":                 "struct{c myChannel}{c: nil}",
		"receive channel":                    "struct{c <-chan int}{c: chan}",
		"send channel":                       "struct{c chan<- int}{c: chan}",
		"function with args":                 "func(int, int) int",
		"custom function with args":          "myFunction",
		"function with multiple return args": "func(int, int) (int, int)",
		"nil function":                       "struct{f func(int) int}{f: nil}",
		"interface type":                     "struct{i interface{foo()}}{i: nil}",
		"map":                                "map[int]int{24: 42}",
		"custom map":                         "myMap{24: 42}",
		"nil map":                            "struct{m map[int]int}{m: nil}",
		"nil custom map":                     "struct{m myMap}{m: nil}",
		"pointer":                            "*struct{}{}",
		"custom pointer":                     "*myStruct{field: nil}",
		"nil pointer":                        "struct{p *int}{p: nil}",
		"nil custom pointer":                 "struct{p myPointer}{p: nil}",
		"list":                               "[]int{1, 2, 3}",
		"custom list":                        "myList{1, 2, 3}",
		"nil list":                           "struct{l []int}{l: nil}",
		"nil custom list":                    "struct{l myList}{l: nil}",
		"long item":                          "[]string{\"foobarbazqux\"}",
		"long subitem":                       "[]struct{foo string}{{foo: \"foobarbazqux\"}}",
		"custom string":                      "myString(\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\")",
		"structure":                          "struct{foo int}{foo: 42}",
		"custom structure":                   "myStruct{field: 42}",
		"unsafe pointer type":                "struct{p Pointer}{p: nil}",
	})
}

func (t tests) expectVerboseTypes() tests {
	return t.expectTypes().expect(map[string]string{
		"false":                              "bool(false)",
		"true":                               "bool(true)",
		"int":                                "int(42)",
		"negative int":                       "int(-42)",
		"array":                              "[3]int{int(1), int(2), int(3)}",
		"custom array":                       "myArray{int(1), int(2), int(3)}",
		"nil channel":                        "struct{c chan int}{c: (chan int)(nil)}",
		"nil custom channel":                 "struct{c myChannel}{c: myChannel(nil)}",
		"receive channel":                    "struct{c <-chan int}{c: <-chan int}",
		"send channel":                       "struct{c chan<- int}{c: chan<- int}",
		"custom function with args":          "myFunction",
		"function with multiple return args": "func(int, int) (int, int)",
		"nil function":                       "struct{f func(int) int}{f: (func(int) int)(nil)}",
		"interface type":                     "struct{i interface{foo()}}{i: interface{foo()}(nil)}",
		"map":                                "map[int]int{int(24): int(42)}",
		"custom map":                         "myMap{int(24): int(42)}",
		"nil map":                            "struct{m map[int]int}{m: (map[int]int)(nil)}",
		"nil custom map":                     "struct{m myMap}{m: myMap(nil)}",
		"custom pointer":                     "*myStruct{field: interface{}(nil)}",
		"nil pointer":                        "struct{p *int}{p: (*int)(nil)}",
		"nil custom pointer":                 "struct{p myPointer}{p: myPointer(nil)}",
		"list":                               "[]int{int(1), int(2), int(3)}",
		"custom list":                        "myList{int(1), int(2), int(3)}",
		"nil list":                           "struct{l []int}{l: ([]int)(nil)}",
		"nil custom list":                    "struct{l myList}{l: myList(nil)}",
		"long item":                          "[]string{string(\"foobarbazqux\")}",
		"long subitem":                       "[]struct{foo string}{struct{foo string}{foo: string(\"foobarbazqux\")}}",
		"string":                             "string(\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\")",
		"structure":                          "struct{foo int}{foo: int(42)}",
		"custom structure":                   "myStruct{field: interface{}(int(42))}",
		"unsafe pointer":                     "Pointer(pointer)",
		"custom unsafe pointer":              "myUnsafePointer(pointer)",
		"unsafe pointer type":                "struct{p Pointer}{p: Pointer(nil)}",
	})
}

func (t tests) expectWrapAll() tests {
	return t.expect(map[string]string{
		"array": `[3]{
	1,
	2,
	3,
}`,
		"custom array": `[3]{
	1,
	2,
	3,
}`,
		"nil channel": `{
	c: nil,
}`,
		"nil custom channel": `{
	c: nil,
}`,
		"receive channel": `{
	c: chan,
}`,
		"send channel": `{
	c: chan,
}`,
		"nil function": `{
	f: nil,
}`,
		"interface type": `{
	i: nil,
}`,
		"map": `map{
	24: 42,
}`,
		"custom map": `map{
	24: 42,
}`,
		"nil map": `{
	m: nil,
}`,
		"nil custom map": `{
	m: nil,
}`,
		"custom pointer": `{
	field: nil,
}`,
		"nil pointer": `{
	p: nil,
}`,
		"nil custom pointer": `{
	p: nil,
}`,
		"list": `[]{
	1,
	2,
	3,
}`,
		"custom list": `[]{
	1,
	2,
	3,
}`,
		"nil list": `{
	l: nil,
}`,
		"nil custom list": `{
	l: nil,
}`,
		"long item": "[]{\n\t\"foobarbazqux\",\n}",
		"long subitem": `[]{
	{
		foo: "foobarbazqux",
	},
}`,
		"structure": `{
	foo: 42,
}`,
		"custom structure": `{
	field: 42,
}`,
		"unsafe pointer type": `{
	p: nil,
}`,
	})
}

func (t tests) expectOnlyLongWrapped() tests {
	return t.expect(map[string]string{
		"long item": `[]{
	"foobarbazqux",
}`,
		"long subitem": `[]{
	{foo: "foobarbazqux"},
}`,
	})
}

func (t tests) expectWrapAllWithTypes() tests {
	return t.expectTypes().expect(map[string]string{
		"array": `[3]int{
	1,
	2,
	3,
}`,
		"custom array": `myArray{
	1,
	2,
	3,
}`,
		"nil channel": `struct{
	c chan int
}{
	c: nil,
}`,
		"nil custom channel": `struct{
	c myChannel
}{
	c: nil,
}`,
		"receive channel": `struct{
	c <-chan int
}{
	c: chan,
}`,
		"send channel": `struct{
	c chan<- int
}{
	c: chan,
}`,
		"function with args": `func(
	int,
	int,
) int`,
		"function with multiple return args": `func(
	int,
	int,
) (
	int,
	int,
)`,
		"nil function": `struct{
	f func(int) int
}{
	f: nil,
}`,
		"interface type": `struct{
	i interface{
		foo()
	}
}{
	i: nil,
}`,
		"map": `map[int]int{
	24: 42,
}`,
		"custom map": `myMap{
	24: 42,
}`,
		"nil map": `struct{
	m map[int]int
}{
	m: nil,
}`,
		"nil custom map": `struct{
	m myMap
}{
	m: nil,
}`,
		"custom pointer": `*myStruct{
	field: nil,
}`,
		"nil pointer": `struct{
	p *int
}{
	p: nil,
}`,
		"nil custom pointer": `struct{
	p myPointer
}{
	p: nil,
}`,
		"list": `[]int{
	1,
	2,
	3,
}`,
		"custom list": `myList{
	1,
	2,
	3,
}`,
		"nil list": `struct{
	l []int
}{
	l: nil,
}`,
		"nil custom list": `struct{
	l myList
}{
	l: nil,
}`,
		"long item": `[]string{
	"foobarbazqux",
}`,
		"long subitem": `[]struct{
	foo string
}{
	{
		foo: "foobarbazqux",
	},
}`,
		"custom string": "myString(\n\t\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\"\n)",
		"structure": `struct{
	foo int
}{
	foo: 42,
}`,
		"custom structure": `myStruct{
	field: 42,
}`,
		"unsafe pointer type": `struct{
	p Pointer
}{
	p: nil,
}`,
	})
}

func (t tests) expectOnlyLongWrappedWithTypes() tests {
	return t.expectTypes().expect(map[string]string{
		"nil channel": `struct{c chan int}{
	c: nil,
}`,
		"nil custom channel": `struct{c myChannel}{
	c: nil,
}`,
		"receive channel": `struct{c <-chan int}{
	c: chan,
}`,
		"send channel": `struct{c chan<- int}{
	c: chan,
}`,
		"nil function": `struct{
	f func(int) int
}{
	f: nil,
}`,
		"interface type": `struct{
	i interface{foo()}
}{
	i: nil,
}`,
		"nil map": `struct{m map[int]int}{
	m: nil,
}`,
		"nil custom map": `struct{m myMap}{
	m: nil,
}`,
		"nil pointer": `struct{p *int}{
	p: nil,
}`,
		"nil custom pointer": `struct{p myPointer}{
	p: nil,
}`,
		"nil list": `struct{l []int}{
	l: nil,
}`,
		"nil custom list": `struct{l myList}{
	l: nil,
}`,
		"long item": `[]string{
	"foobarbazqux",
}`,
		"long subitem": `[]struct{foo string}{
	{foo: "foobarbazqux"},
}`,
		"custom string": "myString(\n\t\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\"\n)",
		"structure": `struct{foo int}{
	foo: 42,
}`,
		"unsafe pointer type": `struct{p Pointer}{
	p: nil,
}`,
	})
}

func (t tests) expectWrapAllWithVerboseTypes() tests {
	return t.expectVerboseTypes().expect(map[string]string{
		"array": `[3]int{
	int(1),
	int(2),
	int(3),
}`,
		"custom array": `myArray{
	int(1),
	int(2),
	int(3),
}`,
		"nil channel": `struct{
	c chan int
}{
	c: (chan int)(nil),
}`,
		"nil custom channel": `struct{
	c myChannel
}{
	c: myChannel(nil),
}`,
		"receive channel": `struct{
	c <-chan int
}{
	c: <-chan int,
}`,
		"send channel": `struct{
	c chan<- int
}{
	c: chan<- int,
}`,
		"function with args": `func(
	int,
	int,
) int`,
		"function with multiple return args": `func(
	int,
	int,
) (
	int,
	int,
)`,
		"nil function": `struct{
	f func(int) int
}{
	f: (func(int) int)(nil),
}`,
		"interface type": `struct{
	i interface{
		foo()
	}
}{
	i: interface{
		foo()
	}(nil),
}`,
		"map": `map[int]int{
	int(24): int(42),
}`,
		"custom map": `myMap{
	int(24): int(42),
}`,
		"nil map": `struct{
	m map[int]int
}{
	m: (map[int]int)(nil),
}`,
		"nil custom map": `struct{
	m myMap
}{
	m: myMap(nil),
}`,
		"custom pointer": `*myStruct{
	field: interface{}(nil),
}`,
		"nil pointer": `struct{
	p *int
}{
	p: (*int)(nil),
}`,
		"nil custom pointer": `struct{
	p myPointer
}{
	p: myPointer(nil),
}`,
		"list": `[]int{
	int(1),
	int(2),
	int(3),
}`,
		"custom list": `myList{
	int(1),
	int(2),
	int(3),
}`,
		"nil list": `struct{
	l []int
}{
	l: ([]int)(nil),
}`,
		"nil custom list": `struct{
	l myList
}{
	l: myList(nil),
}`,
		"long item": `[]string{
	string(
		"foobarbazqux"
	),
}`,
		"long subitem": `[]struct{
	foo string
}{
	struct{
		foo string
	}{
		foo: string(
			"foobarbazqux"
		),
	},
}`,
		"string":        "string(\n\t\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\"\n)",
		"custom string": "myString(\n\t\"\\\\\\\"\\b\\f\\n\\r\\t\\vfoo\"\n)",
		"structure": `struct{
	foo int
}{
	foo: int(42),
}`,
		"custom structure": `myStruct{
	field: interface{}(
		int(42)
	),
}`,
		"unsafe pointer type": `struct{
	p Pointer
}{
	p: Pointer(nil),
}`,
	})
}

func (t tests) expectOnlyLongWrappedWithVerboseTypes() tests {
	return t.expectVerboseTypes().expect(map[string]string{
		"nil channel": `struct{c chan int}{
	c: (chan int)(nil),
}`,
		"nil custom channel": `struct{c myChannel}{
	c: myChannel(nil),
}`,
		"receive channel": `struct{c <-chan int}{
	c: <-chan int,
}`,
		"send channel": `struct{c chan<- int}{
	c: chan<- int,
}`,
		"nil function": `struct{f func(int) int}{
	f: (func(int) int)(nil),
}`,
		"interface type": `struct{i interface{foo()}}{
	i: interface{foo()}(nil),
}`,
		"nil map": `struct{m map[int]int}{
	m: (map[int]int)(nil),
}`,
		"long subitem": `[]struct{foo string}{
	struct{foo string}{
		foo: string("foobarbazqux"),
	},
}`,
		"nil custom pointer": `struct{p myPointer}{
	p: myPointer(nil),
}`,
		"custom structure": `myStruct{
	field: interface{}(int(42)),
}`,
	})
}

func TestSprint(t *testing.T) {
	defaultSet().run(t, Sprint)
}

func TestSprintt(t *testing.T) {
	defaultSet().expectTypes().run(t, Sprintt)
}

func TestSprintv(t *testing.T) {
	defaultSet().expectTypes().expectVerboseTypes().run(t, Sprintv)
}

func TestSprintw(t *testing.T) {
	t.Run("fit everything on a single line", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=8",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet())),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet())),
		)()

		defaultSet().run(t, Sprintw)
	})

	t.Run("wrap everything", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		defaultSet().expectWrapAll().run(t, Sprintw)
	})

	t.Run("wrap long expressions only", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=2",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet())/2),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet())/2),
		)()

		defaultSet().expectOnlyLongWrapped().run(t, Sprintw)
	})

	t.Run("wrap with tolerance", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=2", "LINEWIDTH=9", "LINEWIDTH1=12")()
		tests{{title: "list", value: []int{1, 2, 3}, expect: "[]{1, 2, 3}"}}.run(t, Sprintw)
	})
}

func TestSprintwt(t *testing.T) {
	t.Run("fit everything on a single line", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=8",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet().expectTypes())),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet().expectTypes())),
		)()

		defaultSet().expectTypes().run(t, Sprintwt)
	})

	t.Run("wrap everything", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		defaultSet().expectWrapAllWithTypes().run(t, Sprintwt)
	})

	t.Run("wrap long expressions only", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=2",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet().expectTypes())/2),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet().expectTypes())/2),
		)()

		defaultSet().expectOnlyLongWrappedWithTypes().run(t, Sprintwt)
	})

	t.Run("wrap with tolerance", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=2", "LINEWIDTH=12", "LINEWIDTH1=15")()
		tests{{title: "list", value: []int{1, 2, 3}, expect: "[]int{1, 2, 3}"}}.run(t, Sprintwt)
	})
}

func TestSprintwv(t *testing.T) {
	t.Run("fit everything on a single line", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=8",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet().expectVerboseTypes())),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet().expectVerboseTypes())),
		)()

		defaultSet().expectVerboseTypes().run(t, Sprintwv)
	})

	t.Run("wrap everything", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=0", "LINEWIDTH=0", "LINEWIDTH1=0")()
		defaultSet().expectWrapAllWithVerboseTypes().run(t, Sprintwv)
	})

	t.Run("wrap long expressions only", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=2",
			fmt.Sprintf("LINEWIDTH=%d", expectedMaxWidth(defaultSet().expectVerboseTypes())/2),
			fmt.Sprintf("LINEWIDTH1=%d", expectedMaxWidth(defaultSet().expectVerboseTypes())/2),
		)()

		defaultSet().expectOnlyLongWrappedWithVerboseTypes().run(t, Sprintwv)
	})

	t.Run("wrap with tolerance", func(t *testing.T) {
		defer withEnv(t, "TABWIDTH=2", "LINEWIDTH=27", "LINEWIDTH1=30")()
		tests{{
			title:  "list",
			value:  []int{1, 2, 3},
			expect: "[]int{int(1), int(2), int(3)}",
		}}.run(t, Sprintwv)
	})
}

func TestNoLonger(t *testing.T) {
	const expect = `[]{"foobarbaz"}`
	o := []string{"foobarbaz"}
	defer withEnv(t, "TABWIDTH=8", "LINEWIDTH=9", "LINEWIDTH1=9")()
	s := Sprintw(o)
	if s != expect {
		t.Fatalf("expected: %s, got: %s", expect, s)
	}
}

func TestSprintMultipleObjects(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		const expect = "{fooBarBaz: 42} {fooBarBaz: 42}"
		o := struct{ fooBarBaz int }{42}
		s := Sprint(o, o)
		if s != expect {
			t.Fatalf("expected: %s, got: %s", expect, s)
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		const expect = "{fooBarBaz: 42}\n{fooBarBaz: 42}"
		o := struct{ fooBarBaz int }{42}
		s := Sprintw(o, o)
		if s != expect {
			t.Fatalf("expected: %s, got: %s", expect, s)
		}
	})
}

func TestInvalidEasterEgg(t *testing.T) {
	defer withEnv(t, "TABWIDTH=foo")()
	const expect = "{fooBarBaz: 42}"
	o := struct{ fooBarBaz int }{42}
	s := Sprintw(o)
	if s != expect {
		t.Fatalf("expected: %s, got: %s", expect, s)
	}
}

func TestAnonymousFields(t *testing.T) {
	type Foo struct{ Bar int }
	type Baz struct{ Foo }
	o := Baz{Foo: Foo{Bar: 42}}
	t.Run("without types", func(t *testing.T) {
		const expect = `{Foo: {Bar: 42}}`
		s := Sprint(o)
		if s != expect {
			t.Fatalf("expected: %s, got: %s", expect, s)
		}
	})

	t.Run("with types", func(t *testing.T) {
		const expect = `Baz{Foo: {Bar: 42}}`
		s := Sprintt(o)
		if s != expect {
			t.Fatalf("expected: %s, got: %s", expect, s)
		}
	})

	t.Run("with verbose types", func(t *testing.T) {
		const expect = `Baz{Foo: Foo{Bar: int(42)}}`
		s := Sprintv(o)
		if s != expect {
			t.Fatalf("expected: %s, got: %s", expect, s)
		}
	})
}

func TestVariadicFunction(t *testing.T) {
	const expect = `func(...int)`
	f := func(...int) {}
	s := Sprintt(f)
	if s != expect {
		t.Fatalf("expected: %s, got: %s", expect, s)
	}
}

func TestSortedMap(t *testing.T) {
	const testCount = 9
	m := map[int]int{1: 2, 3: 4, 5: 6}
	t.Run("sorted", func(t *testing.T) {
		const expect = `map{1: 2, 3: 4, 5: 6}`
		for i := 0; i < testCount; i++ {
			s := Sprint(m)
			if s != expect {
				t.Fatalf("expected: %s, got: %s", expect, s)
			}
		}
	})

	t.Run("random", func(t *testing.T) {
		defer withEnv(t, "MAPSORT=0")()
		for i := 0; i < testCount; i++ {
			s := Sprint(m)
			for _, entry := range []string{"1: 2", "3: 4", "5: 6"} {
				if !strings.Contains(s, entry) {
					t.Fatalf("missing entry: %s", entry)
				}
			}
		}
	})
}

func TestBytes(t *testing.T) {
	const expectNotWrapped = `[]{00 01 02 03 04 05 06 07 08 09 0a 0b}`
	const expectWrapped = `[]{
	00 01 02 03 04 05
	06 07 08 09 0a 0b
}`

	b := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	t.Run("not wrapped", func(t *testing.T) {
		s := Sprint(b)
		if s != expectNotWrapped {
			t.Fatalf("expected: %s, got: %s", expectNotWrapped, s)
		}
	})

	t.Run("wrapped", func(t *testing.T) {
		defer withEnv(
			t,
			"TABWIDTH=2",
			fmt.Sprintf("LINEWIDTH=%d", len(expectNotWrapped)/2),
			fmt.Sprintf("LINEWIDTH1=%d", len(expectNotWrapped)/2),
		)()

		s := Sprintw(b)
		if s != expectWrapped {
			t.Fatalf("expected: %s, got: %s", expectWrapped, s)
		}
	})
}

func TestNonWrapperNodes(t *testing.T) {
	const expect = `map[struct{foo int; bar int}]struct{
	foo int
	bar int
}{}`

	defer withEnv(t, "TABWIDTH=2", "LINEWIDTH=27", "LINEWIDTH1=30")()
	o := map[struct{foo int; bar int}]struct{foo int; bar int}{}
	s := Sprintwt(o)
	if s != expect {
		t.Fatalf("expected: %s, got: %s", expect, s)
	}
}
