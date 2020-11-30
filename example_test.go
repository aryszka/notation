package notation_test

import (
	"os"

	"github.com/aryszka/notation"
)

type bike struct {
	frame      frame
	driveTrain driveTrain
	wheels     []wheel
	handlebar  handlebar
	saddle     saddle
}
type frame struct {
	fork            fork
	saddlePost      saddlePost
	bottomBracket   *bracket
	frontDerailleur *derailleur
	rearDerailleur  *derailleur
	rearBrake       *brake
	rearWheel       *wheel
}
type driveTrain struct {
	bottomBracket bracket
	crank         crank
	brakes        []brake
	derailleurs   []derailleur
	cassette      cassette
	chain         chain
	levers        []lever
}
type wheel struct {
	size     float64
	cassette *cassette
}
type handlebar struct{ levers []*lever }
type saddle struct{}
type fork struct {
	wheel      *wheel
	handlebar  *handlebar
	frontBrake *brake
}
type saddlePost struct{ saddle *saddle }
type bracket struct{ crank *crank }
type derailleur struct{ gears int }
type brake struct{ discSize float64 }
type crank struct {
	wheels int
	chain  *chain
}
type cassette struct {
	wheels int
	chain  *chain
}
type chain struct{}
type lever struct{ withShift bool }

func Example() {
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
		wheels:    []wheel{{size: 700}, {size: 700}},
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

	notation.Fprintw(os.Stdout, b)

	// Output:
	//
	// {
	// 	frame: {
	// 		fork: {
	// 			wheel: {size: 700, cassette: nil},
	// 			handlebar: {
	// 				levers: []{
	// 					{withShift: true},
	// 					{withShift: true},
	// 				},
	// 			},
	// 			frontBrake: {discSize: 160},
	// 		},
	// 		saddlePost: {saddle: {}},
	// 		bottomBracket: {crank: {wheels: 2, chain: {}}},
	// 		frontDerailleur: {gears: 2},
	// 		rearDerailleur: {gears: 11},
	// 		rearBrake: {discSize: 140},
	// 		rearWheel: {
	// 			size: 700,
	// 			cassette: {wheels: 11, chain: {}},
	// 		},
	// 	},
	// 	driveTrain: {
	// 		bottomBracket: {crank: {wheels: 2, chain: {}}},
	// 		crank: {wheels: 2, chain: {}},
	// 		brakes: []{{discSize: 160}, {discSize: 140}},
	// 		derailleurs: []{{gears: 2}, {gears: 11}},
	// 		cassette: {wheels: 11, chain: {}},
	// 		chain: {},
	// 		levers: []{{withShift: true}, {withShift: true}},
	// 	},
	// 	wheels: []{
	// 		{size: 700, cassette: nil},
	// 		{size: 700, cassette: {wheels: 11, chain: {}}},
	// 	},
	// 	handlebar: {levers: []{{withShift: true}, {withShift: true}}},
	// 	saddle: {},
	// }
}

func Example_int() {
	i := 42
	notation.Fprintt(os.Stdout, i)

	// Output:
	//
	// 42
}

func Example_short_string() {
	s := `foobar
baz`
	notation.Fprintw(os.Stdout, s)

	// Output:
	//
	// "foobar\nbaz"
}

func Example_long_string() {
	s := `The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog.`

	notation.Fprintw(os.Stdout, s)

	// Output:
	//
	// `The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
	// quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
	// fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
	// over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
	// dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
	// quick brown fox jumps over the lazy dog.`
}

func Example_slice() {
	l := []int{1, 2, 3}
	notation.Fprint(os.Stdout, l)

	// Output:
	//
	// []{1, 2, 3}
}

func Example_array() {
	a := [...]int{1, 2, 3}
	notation.Fprint(os.Stdout, a)

	// Output:
	//
	// [3]{1, 2, 3}
}

func Example_bytes() {
	b := []byte(
		`The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown
fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps
over the lazy dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy
dog. The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog. The
quick brown fox jumps over the lazy dog.`,
	)

	notation.Fprintwt(os.Stdout, b)

	// Output:
	//
	// []byte{
	// 	54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a
	// 	75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f
	// 	67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f
	// 	78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79
	// 	20 64 6f 67 2e 20 54 68 65 0a 71 75 69 63 6b 20 62 72 6f 77 6e
	// 	20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c
	// 	61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72
	// 	6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68
	// 	65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b
	// 	20 62 72 6f 77 6e 0a 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72
	// 	20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75
	// 	69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f
	// 	76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65
	// 	20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70
	// 	73 0a 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e 20
	// 	54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f 78 20 6a
	// 	75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79 20 64 6f
	// 	67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e 20 66 6f
	// 	78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c 61 7a 79
	// 	0a 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72 6f 77 6e
	// 	20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68 65 20 6c
	// 	61 7a 79 20 64 6f 67 2e 20 54 68 65 20 71 75 69 63 6b 20 62 72
	// 	6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72 20 74 68
	// 	65 20 6c 61 7a 79 20 64 6f 67 2e 20 54 68 65 0a 71 75 69 63 6b
	// 	20 62 72 6f 77 6e 20 66 6f 78 20 6a 75 6d 70 73 20 6f 76 65 72
	// 	20 74 68 65 20 6c 61 7a 79 20 64 6f 67 2e
	// }
}

func Example_maps_sorted_by_keys() {
	m := map[string]int{"b": 1, "c": 2, "a": 3}
	notation.Fprint(os.Stdout, m)

	// Output:
	//
	// map{"a": 3, "b": 1, "c": 2}
}

func Example_function() {
	f := func(int) int { return 42 }
	notation.Fprint(os.Stdout, f)

	// Output:
	//
	// func()
}

func Example_function_signature() {
	f := func(int) int { return 42 }
	notation.Fprintt(os.Stdout, f)

	// Output:
	//
	// func(int) int
}

func Example_named_type() {
	type t struct{ foo int }
	v := t{42}
	notation.Fprintt(os.Stdout, v)

	// Output:
	//
	// t{foo: 42}
}

func Example_unnamed() {
	v := struct{ foo int }{42}
	notation.Fprintt(os.Stdout, v)

	// Output:
	//
	// struct{foo int}{foo: 42}
}

func Example_type_inferred() {
	v := []struct{ foo int }{{42}, {84}}
	notation.Fprintt(os.Stdout, v)

	// Output:
	//
	// []struct{foo int}{{foo: 42}, {foo: 84}}
}

func Example_cyclic_reference() {
	l := []interface{}{"foo"}
	l[0] = l
	notation.Fprint(os.Stdout, l)

	// Output:
	//
	// r0=[]{r0}
}
