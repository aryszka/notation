package notation

import (
	"reflect"
	"strings"
)

type opts int

const none opts = 0

const (
	wrap opts = 1 << iota
	types
	skipTypes
	allTypes
)

func sprintValues(o opts, v []interface{}) string {
	s := make([]string, len(v))
	for i := range v {
		if v[i] == nil {
			s[i] = "nil"
			continue
		}

		s[i] = sprint(o, reflect.ValueOf(v[i]))
	}

	sep := " "
	if o&wrap != 0 {
		sep = "\n"
	}

	return strings.Join(s, sep)
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
