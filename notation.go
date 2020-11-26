/*
Package notation can be used to print (or sprint) Go objects with optional wrapping (and indentation) and
optional type information.
*/
package notation

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type opts int

const none opts = 0

const (
	wrap opts = 1 << iota
	types
	skipTypes
	allTypes
	randomMaps
)

type wrapLen struct {
	first, max, last int
}

type nodeRef struct {
	id, refCount int
}

type pending struct {
	values    map[uintptr]nodeRef
	idCounter int
}

type node struct {
	parts    []interface{}
	len      int
	wrapLen  wrapLen
	fullWrap wrapLen
	wrap     bool
}

type str struct {
	val    string
	raw    string
	rawLen wrapLen
	useRaw bool
}

type wrapMode int

const (
	block wrapMode = iota
	line
)

type wrapper struct {
	mode        wrapMode
	sep, suffix string
	items       []node
	lineEnds    []int
}

type writer struct {
	w   io.Writer
	n   int
	err error
}

var stderr io.Writer = os.Stderr

func nodeOf(parts ...interface{}) node {
	return node{parts: parts}
}

// used only for debugging
func (n node) String() string {
	var b bytes.Buffer
	w := &writer{w: &b}
	fprint(w, 0, n)
	return b.String()
}

func (s str) String() string {
	if s.useRaw {
		return s.raw
	}

	return s.val
}

func (w *writer) write(o interface{}) {
	if w.err != nil {
		return
	}

	n, err := fmt.Fprint(w.w, o)
	w.n += n
	w.err = err
}

func (w *writer) blankLine() {
	w.write("\n")
}

func (w *writer) tabs(n int) {
	for i := 0; i < n; i++ {
		w.write("\t")
	}
}

func (w *writer) line(t int) {
	w.blankLine()
	w.tabs(t)
}

func config(name string, dflt int) int {
	s := os.Getenv(name)
	if s == "" {
		s = os.Getenv(strings.ToLower(name))
	}

	if s == "" {
		return dflt
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return dflt
	}

	return v
}

func fprintValues(w io.Writer, o opts, v []interface{}) (int, error) {
	tab := config("TABWIDTH", 8)
	cols0 := config("LINEWIDTH", 80-tab)
	cols1 := config("LINEWIDTH1", (cols0+tab)*3/2-tab)
	sortMaps := config("MAPSORT", 1)
	if sortMaps == 0 {
		o |= randomMaps
	}

	wr := &writer{w: w}
	for i, vi := range v {
		if wr.err != nil {
			return wr.n, wr.err
		}

		if i > 0 {
			if o&wrap == 0 {
				wr.write(" ")
			} else {
				wr.write("\n")
			}
		}

		if vi == nil {
			fprint(wr, 0, nodeOf("nil"))
			continue
		}

		p := &pending{values: make(map[uintptr]nodeRef)}
		n := reflectValue(o, p, reflect.ValueOf(vi))
		if o&wrap != 0 {
			n = nodeLen(tab, n)
			n = wrapNode(tab, cols0, cols0, cols1, n)
		}

		fprint(wr, 0, n)
	}

	return wr.n, wr.err
}

func printValues(o opts, v []interface{}) (int, error) {
	return fprintValues(stderr, o, v)
}

func printlnValues(o opts, v []interface{}) (int, error) {
	n, err := fprintValues(stderr, o, v)
	if err != nil {
		return n, err
	}

	nn, err := stderr.Write([]byte("\n"))
	return n + nn, err
}

func sprintValues(o opts, v []interface{}) string {
	var b bytes.Buffer
	fprintValues(&b, o, v)
	return b.String()
}

// Fprint prints the provided objects to the provided writer. When multiple objects are printed, they'll be
// separated by a space.
func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, none, v)
}

// Fprintw prints the provided objects to the provided writer, with wrapping (and indentation) where necessary.
// When multiple objects are printed, they'll be separated by a newline.
func Fprintw(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap, v)
}

// Fprintt prints the provided objects to the provided writer with moderate type information. When multiple
// objects are printed, they'll be separated by a space.
func Fprintt(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, types, v)
}

// Fprintwt prints the provided objects to the provided writer, with wrapping (and indentation) where necessary,
// and with moderate type information. When multiple objects are printed, they'll be separated by a newline.
func Fprintwt(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap|types, v)
}

// Fprintv prints the provided objects to the provided writer with verbose type information. When multiple
// objects are printed, they'll be separated by a space.
func Fprintv(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, allTypes, v)
}

// Fprintwv prints the provided objects to the provided writer, with wrapping (and indentation) where necessary,
// and with verbose type information. When multiple objects are printed, they'll be separated by a newline.
func Fprintwv(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap|allTypes, v)
}

// Print prints the provided objects to stderr. When multiple objects are printed, they'll be separated by a
// space.
func Print(v ...interface{}) (int, error) {
	return printValues(none, v)
}

// Printw prints the provided objects to stderr, with wrapping (and indentation) where necessary. When multiple
// objects are printed, they'll be separated by a newline.
func Printw(v ...interface{}) (int, error) {
	return printValues(wrap, v)
}

// Printt prints the provided objects to stderr with moderate type information. When multiple objects are
// printed, they'll be separated by a space.
func Printt(v ...interface{}) (int, error) {
	return printValues(types, v)
}

// Printwt prints the provided objects to stderr, with wrapping (and indentation) where necessary, and with
// moderate type information. When multiple objects are printed, they'll be separated by a newline.
func Printwt(v ...interface{}) (int, error) {
	return printValues(wrap|types, v)
}

// Printv prints the provided objects to stderr with verbose type information. When multiple objects are
// printed, they'll be separated by a space.
func Printv(v ...interface{}) (int, error) {
	return printValues(allTypes, v)
}

// Printwv prints the provided objects to stderr, with wrapping (and indentation) where necessary, and with
// verbose type information. When multiple objects are printed, they'll be separated by a newline.
func Printwv(v ...interface{}) (int, error) {
	return printValues(wrap|allTypes, v)
}

// Println prints the provided objects to stderr with a closing newline. When multiple objects are printed,
// they'll be separated by a space.
func Println(v ...interface{}) (int, error) {
	return printlnValues(none, v)
}

// Printlnw prints the provided objects to stderr with a closing newline, with wrapping (and indentation) where
// necessary. When multiple objects are printed, they'll be separated by a newline.
func Printlnw(v ...interface{}) (int, error) {
	return printlnValues(wrap, v)
}

// Printlnt prints the provided objects to stderr with a closing newline, and with moderate type information. When
// multiple objects are printed, they'll be separated by a space.
func Printlnt(v ...interface{}) (int, error) {
	return printlnValues(types, v)
}

// Printlnwt prints the provided objects to stderr with a closing newline, with wrapping (and indentation) where
// necessary, and with moderate type information. When multiple objects are printed, they'll be separated by a
// newline.
func Printlnwt(v ...interface{}) (int, error) {
	return printlnValues(wrap|types, v)
}

// Printlnv prints the provided objects to stderr with a closing newline, and with verbose type information. When
// multiple objects are printed, they'll be separated by a space.
func Printlnv(v ...interface{}) (int, error) {
	return printlnValues(allTypes, v)
}

// Printlnwv prints the provided objects to stderr with a closing newline, with wrapping (and indentation) where
// necessary, and with verbose type information. When multiple objects are printed, they'll be separated by a
// newline.
func Printlnwv(v ...interface{}) (int, error) {
	return printlnValues(wrap|allTypes, v)
}

// Sprint returns the string representation of the Go objects. When multiple objects are provided, they'll be
// seprated by a space.
func Sprint(v ...interface{}) string {
	return sprintValues(none, v)
}

// Sprintw returns the string representation of the Go objects, with wrapping (and indentation) where necessary.
// When multiple objects are provided, they'll be seprated by a newline.
func Sprintw(v ...interface{}) string {
	return sprintValues(wrap, v)
}

// Sprintt returns the string representation of the Go objects, with moderate type information. When multiple
// objects are provided, they'll be seprated by a space.
func Sprintt(v ...interface{}) string {
	return sprintValues(types, v)
}

// Sprintwt returns the string representation of the Go objects, with wrapping (and indentation) where necessary,
// and with moderate type information. When multiple objects are provided, they'll be seprated by a newline.
func Sprintwt(v ...interface{}) string {
	return sprintValues(wrap|types, v)
}

// Sprintv returns the string representation of the Go objects, with verbose type information. When multiple
// objects are provided, they'll be seprated by a space.
func Sprintv(v ...interface{}) string {
	return sprintValues(allTypes, v)
}

// Sprintwv returns the string representation of the Go objects, with wrapping (and indentation) where necessary,
// and with verbose type information. When multiple objects are provided, they'll be seprated by a newline.
func Sprintwv(v ...interface{}) string {
	return sprintValues(wrap|allTypes, v)
}
