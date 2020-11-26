# Notation - print Go objects

[![codecov](https://codecov.io/gh/aryszka/notation/branch/master/graph/badge.svg?token=7M18MEAVQW)](https://codecov.io/gh/aryszka/notation)

This package can be used to print (or sprint) Go objects for debugging purposes, with optional wrapping
(indentation) and optional type information.

### Alternatives

Notation is similar to the following great, more established and mature packages:

- [go-spew](https://github.com/davecgh/go-spew/)
- [litter](https://github.com/sanity-io/litter)
- [utter](https://github.com/kortschak/utter)

Notation differs from these primarily in the 'flavor' of printing and the package interface.

### Installation

`go get github.com/aryszka/notation`

### Usage

Pass in the Go object(s) to be printed to one of the notation functions. These functions can be categorized into
four groups:

- **Println:** the Println type functions print Go objects to stderr with an additional newline at the end.
- **Print:** similar to Println but without the extra newline.
- **Fprint:** the Fprint type functions print Go objects to an arbitrary writer passed in as an argument.
- **Sprint:** the Sprint type functions return the string representation of Go objects instead of printing them.

The format and the verbosity can be controlled with the suffixed variants of the above functions. By default,
the input arguments are printed without types on a single line. The following suffixes are available:

- **`w`**: wrap the output based on Go-style indentation
- **`t`**: print moderately verbose type information, omitting where it can be trivially inferred
- **`v`**: print verbose type information

When **`t`** or **`v`** are used together with **`w`**, they must follow **`w`**. The suffixes **`t`** and
**`v`** cannot be used together.

Wrapping is **not eager**. It doesn't wrap a line when it can fit on 72 columns. It also tolerates longer lines up
to 112 columns, when the output is considered to be more readable that way. This means very simple Go objects
are not wrapped even with the **`w`** variant of the functions.

For the available functions, see also the [godoc](https://godoc.org/github.com/aryszka/notation).
(Alternatively: [pkg.go.dev](https://pkg.go.dev/github.com/aryszka/notation).)

### Example

Assuming to have the required types defined, if we do the following:

```
b := bike{
	frame: frame{
		fork:       fork{},
		saddlePost: saddlePost{},
	},
	driveTrain: driveTrain{
		bottomBracket: bracket{},
		crank:         crank{wheels: 2},
		brakes:        []brake{{discSize: 160}, {discSize: 140}},
		derailleurs:   []derailleur{{gears: 2}, {gears: 11}},
		cassette:      cassette{wheels: 11},
		chain:         chain{},
		levers:        []lever{{true}, {true}},
	},
	wheels:    []wheel{{size: 70}, {size: 70}},
	handlebar: handlebar{},
	saddle:    saddle{},
}

b.frame.fork.wheel = &b.wheels[0]
b.frame.fork.handlebar = &b.handlebar
b.frame.fork.handlebar.levers = []*lever{&b.driveTrain.levers[0], &b.driveTrain.levers[1]}
b.frame.fork.frontBrake = &b.driveTrain.brakes[0]
b.frame.saddlePost.saddle = &b.saddle
b.frame.bottomBracket = &b.driveTrain.bottomBracket
b.frame.frontDerailleur = &b.driveTrain.derailleurs[0]
b.frame.rearDerailleur = &b.driveTrain.derailleurs[1]
b.frame.rearBrake = &b.driveTrain.brakes[1]
b.frame.rearWheel = &b.wheels[1]
b.frame.bottomBracket.crank = &b.driveTrain.crank
b.frame.bottomBracket.crank.chain = &b.driveTrain.chain
b.frame.rearWheel.cassette = &b.driveTrain.cassette
b.frame.rearWheel.cassette.chain = &b.driveTrain.chain

s := notation.Sprintw(b)
```

We get the following string:

```
{
	frame: {
		fork: {
			wheel: {size: 70, cassette: nil},
			handlebar: {
				levers: []{
					{withShift: true},
					{withShift: true},
				},
			},
			frontBrake: {discSize: 160},
		},
		saddlePost: {saddle: {}},
		bottomBracket: {crank: {wheels: 2, chain: {}}},
		frontDerailleur: {gears: 2},
		rearDerailleur: {gears: 11},
		rearBrake: {discSize: 140},
		rearWheel: {size: 70, cassette: {wheels: 11, chain: {}}},
	},
	driveTrain: {
		bottomBracket: {crank: {wheels: 2, chain: {}}},
		crank: {wheels: 2, chain: {}},
		brakes: []{{discSize: 160}, {discSize: 140}},
		derailleurs: []{{gears: 2}, {gears: 11}},
		cassette: {wheels: 11, chain: {}},
		chain: {},
		levers: []{{withShift: true}, {withShift: true}},
	},
	wheels: []{
		{size: 70, cassette: nil},
		{size: 70, cassette: {wheels: 11, chain: {}}},
	},
	handlebar: {levers: []{{withShift: true}, {withShift: true}}},
	saddle: {},
}
```

Using `notation.Sprintwv` instead of `notation.Sprintw`, we would get the following string:

```
bike{
	frame: frame{
		fork: fork{
			wheel: *wheel{
				size: float64(70),
				cassette: (*cassette)(nil),
			},
			handlebar: *handlebar{
				levers: []*lever{
					*lever{withShift: bool(true)},
					*lever{withShift: bool(true)},
				},
			},
			frontBrake: *brake{discSize: float64(160)},
		},
		saddlePost: saddlePost{saddle: *saddle{}},
		bottomBracket: *bracket{
			crank: *crank{
				wheels: int(2),
				chain: *chain{},
			},
		},
		frontDerailleur: *derailleur{gears: int(2)},
		rearDerailleur: *derailleur{gears: int(11)},
		rearBrake: *brake{discSize: float64(140)},
		rearWheel: *wheel{
			size: float64(70),
			cassette: *cassette{
				wheels: int(11),
				chain: *chain{},
			},
		},
	},
	driveTrain: driveTrain{
		bottomBracket: bracket{
			crank: *crank{
				wheels: int(2),
				chain: *chain{},
			},
		},
		crank: crank{wheels: int(2), chain: *chain{}},
		brakes: []brake{
			brake{discSize: float64(160)},
			brake{discSize: float64(140)},
		},
		derailleurs: []derailleur{
			derailleur{
				gears: int(2),
			},
			derailleur{
				gears: int(11),
			},
		},
		cassette: cassette{wheels: int(11), chain: *chain{}},
		chain: chain{},
		levers: []lever{
			lever{withShift: bool(true)},
			lever{withShift: bool(true)},
		},
	},
	wheels: []wheel{
		wheel{size: float64(70), cassette: (*cassette)(nil)},
		wheel{
			size: float64(70),
			cassette: *cassette{wheels: int(11), chain: *chain{}},
		},
	},
	handlebar: handlebar{
		levers: []*lever{
			*lever{withShift: bool(true)},
			*lever{withShift: bool(true)},
		},
	},
	saddle: saddle{},
}
```

### Subtleties

Notation doesn't provide configuration options, we can just pick the preferred function and call it with the
objects to be printed. The following details describe some of the behavior to be expected in case of various
input objects.

##### Numbers

Numbers are printed based on the `fmt` package's default formatting. When printing with moderate type
information, the type for the `int`, default width signed integers, will be omitted.

```
i := 42
notation.Printlnt(i)
```

Output:

```
42
```

##### Strings

When printing strings, by default they are escaped using the `strconv.Quote` function. However, when wrapping
long strings, and the string contains a newline and doesn't contain a backquote, then the string is printed
as a raw string literal, delimited by backquotes.

Short string:

```
s := `foobar
baz`
notation.Printlnw(s)
```

Output:

```
"foobar\nbaz"
```

Long string:

```
s := `The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog.`

notation.Println(s)
```

Output:

```
`The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog.`
```

##### Arrays/Slices

Slices are are printed by printing their elements between braces, prefixed either by '[]' or the type of the
slice. Example:

```
l := []int{1, 2, 3}
notation.Println(l)
```

Output:

```
[]{1, 2, 3}
```

To differentiate arrays from slices, arrays are always prefixed with their type or square brackets containing
the length of the array:

```
a := [...]{1, 2, 3}
notation.Println(a)
```

Output:

```
[3]{1, 2, 3}
```

##### Bytes

When the type of a slice is `uint8`, or an alias of it, e.g. `byte`, then it is printed as []byte, with the hexa
representation of its bytes:

```
b := []byte(
	`The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog.`,
)

notation.Printlnwt(b)
```

Output:

```
[]byte{
	54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a
	75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f
	67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f
	78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79
	20 64 6f 67 2e 20 54 68 65 0a 71 75 69 63 6b 20 62 72 6f 77 6e
	20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c
	61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72
	6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68
	65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b
	20 62 72 6f 77 6e 0a 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72
	20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75
	69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f
	76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65
	20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70
	73 0a 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20
	54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a
	75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f
	67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f
	78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79
	0a 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e
	20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c
	61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72
	6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68
	65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 0a 71 75 69 63 6b
	20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72
	20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e
}
```

##### Maps

Maps are printed with their entries sorted by the string representation of their keys:

```
m := map[string]int{"b": 1, "c": 2, "a": 3}
notation.Println(m)
```

Output:

```
map{"a": 3, "b": 1, "c": 2}
```

This way a map is printed always the same way. If, for a reason, this is undesired, then this behavior can be
disabled via the `MAPSORT=0` environment variable.

##### Hidden values: channels, functions

Certain values, like channels and functions are printed without expanding their internals, e.g. channel state or
function body. When printing with types, the signature of these objects is printed:

```
f := func(int) int { return 42 }
notation.Println(f)
```

Output:

```
func()
```

With types:

```
f := func(int) int { return 42 }
notation.Printlnt(f)
```

Output:

```
func(int) int
```

##### Wrapping

The 'w' variant of the printing functions wraps the output with Go style indentation where the lines would be
too long otherwise. The wrapping is not eager, it only aims for fitting the lines on 72 columns. To measure the
indentation, it assumes 8 character width tabs. In certain cases, it tolerates longer lines up to 112 columns,
when the output would probably more readable that way. Of course, readability is subjective.

As a hidden feature, when it's really necessary, it is possible to change the above control values via
environment variables. TABWIDTH controls the measuring of the indentation. LINEWIDTH sets the aimed column width
of the printed lines. LINEWIDTH1 sets the tolerated threshold for those lines that are allowed to exceed the
default line width. E.g. if somebody uses two-character wide tabs in their console, they can use the package
like this:

```
TABWIDTH=2 go test -v -count 1
```

As a consequence, it is also possible to forcibly wrap all lines:

```
TABWIDTH=0 LINEWIDTH=0 LINEWIDTH1=0 go test -v -count 1
```

##### Types

Using the 't' or 'v' suffixed variants of the printing functions, notation prints the types together with the
values. When the name of a type is available, the name is printed instead of the literal representation of the
type. The package path is not printed.

Named type:

```
type t struct{foo int}
v := t{42}
notation.Printlnt(v)
```

Output:

```
t{foo: 42}
```

Unnamed type:

```
v := struct{foo int}{42}
notation.Printlnt(v)
```

Output:

```
struct{foo int}{foo: 42}
```

Using the 't' suffixed variants of the printing functions, displaying only moderately verbose type information,
the types of certain values is omitted, where it can be inferred from the context:

```
v := []struct{foo int}{{42}, {84}}
notation.Printlnt(os.Stdout, v)
```

Output:

```
[]struct{foo int}{{foo: 42}, {foo: 84}}
```

##### Cyclic references

Cyclic references are detected based on an approach similar to the one in the stdlib's reflect.DeepEqual
function. Such occurrences are displayed in the output with references:

```
l := []interface{}{"foo"}
l[0] = l
notation.Fprint(os.Stdout, l)
```

Output:

```
r0=[]{r0}
```

