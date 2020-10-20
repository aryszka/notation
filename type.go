package notation

import (
	"fmt"
	"reflect"
	"strings"
)

func funcBase(t reflect.Type) string {
	args := func(num func() int, typ func(int) reflect.Type) []string {
		t := make([]string, num())
		for i := 0; i < num(); i++ {
			t[i] = sprintType(typ(i))
		}

		return t
	}

	in := args(t.NumIn, t.In)
	out := args(t.NumOut, t.Out)

	var outs string
	if len(out) == 1 {
		outs = out[0]
	} else if len(out) > 1 {
		outs = fmt.Sprintf("(%s)", strings.Join(out, ", "))
	}

	var s string
	if outs == "" {
		s = fmt.Sprintf("(%s)", strings.Join(in, ", "))
	} else {
		s = fmt.Sprintf("(%s) %s", strings.Join(in, ", "), outs)
	}

	return s
}

func arrayType(t reflect.Type) string {
	return fmt.Sprintf("[%d]%s", t.Len(), sprintType(t.Elem()))
}

func chanType(t reflect.Type) string {
	var prefix string
	switch t.ChanDir() {
	case reflect.RecvDir:
		prefix = "<-chan"
	case reflect.SendDir:
		prefix = "chan<-"
	default:
		prefix = "chan"
	}

	return fmt.Sprintf("%s %s", prefix, sprintType(t.Elem()))
}

func funcType(t reflect.Type) string {
	return fmt.Sprintf("func%s", funcBase(t))
}

func interfaceType(t reflect.Type) string {
	var m []string
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		m = append(m, fmt.Sprintf("%s%s", method.Name, funcBase(method.Type)))
	}

	return fmt.Sprintf("interface{%s}", strings.Join(m, "; "))
}

func mapType(t reflect.Type) string {
	return fmt.Sprintf("map[%s]%s", sprintType(t.Key()), sprintType(t.Elem()))
}

func pointerType(t reflect.Type) string {
	return fmt.Sprintf("*%s", sprintType(t.Elem()))
}

func listType(t reflect.Type) string {
	return fmt.Sprintf("[]%s", sprintType(t.Elem()))
}

func structType(t reflect.Type) string {
	f := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		f[i] = fmt.Sprintf("%s %s", fi.Name, sprintType(fi.Type))
	}

	return fmt.Sprintf("struct{%s}", strings.Join(f, "; "))
}

func sprintType(t reflect.Type) string {
	if t.Name() != "" {
		return t.Name()
	}

	switch t.Kind() {
	case reflect.Array:
		return arrayType(t)
	case reflect.Chan:
		return chanType(t)
	case reflect.Func:
		return funcType(t)
	case reflect.Interface:
		return interfaceType(t)
	case reflect.Map:
		return mapType(t)
	case reflect.Ptr:
		return pointerType(t)
	case reflect.Slice:
		return listType(t)
	case reflect.Struct:
		return structType(t)
	default:
		return "<invalid>"
	}
}
