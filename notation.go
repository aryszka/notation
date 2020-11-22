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

type valueKey struct {
	typ reflect.Type
	ptr uintptr
}

type nodeRef struct {
	id, refCount int
}

type pending struct {
	values    map[valueKey]nodeRef
	idCounter int
}

type node struct {
	len      int
	wrapLen  wrapLen
	fullWrap wrapLen
	wrap     bool
	parts    []interface{}
}

type str struct {
	val    string
	raw    string
	useRaw bool
	rawLen wrapLen
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

func nodeOf(parts ...interface{}) node {
	return node{parts: parts}
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

		n := reflectValue(
			o,
			&pending{values: make(map[valueKey]nodeRef)},
			reflect.ValueOf(vi),
		)

		if o&wrap != 0 {
			n = nodeLen(tab, n)
			n = wrapNode(tab, cols0, cols0, cols1, n)
		}

		fprint(wr, 0, n)
	}

	return wr.n, wr.err
}

func sprintValues(o opts, v []interface{}) string {
	var b bytes.Buffer
	fprintValues(&b, o, v)
	return b.String()
}

func Fprint(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, none, v)
}

func Fprintw(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap, v)
}

func Fprintt(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, types, v)
}

func Fprintwt(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap|types, v)
}

func Fprintv(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, allTypes, v)
}

func Fprintwv(w io.Writer, v ...interface{}) (int, error) {
	return fprintValues(w, wrap|allTypes, v)
}

func Sprint(v ...interface{}) string {
	return sprintValues(none, v)
}

func Sprintw(v ...interface{}) string {
	return sprintValues(wrap, v)
}

func Sprintt(v ...interface{}) string {
	return sprintValues(types, v)
}

func Sprintwt(v ...interface{}) string {
	return sprintValues(wrap|types, v)
}

func Sprintv(v ...interface{}) string {
	return sprintValues(allTypes, v)
}

func Sprintwv(v ...interface{}) string {
	return sprintValues(wrap|allTypes, v)
}
