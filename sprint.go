package notation

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func withType(o opts) (opts, bool, bool) {
	if o&types == 0 && o&allTypes == 0 {
		return o, false, false
	}

	if o&skipTypes != 0 && o&allTypes == 0 {
		return o &^ skipTypes, false, false
	}

	return o, true, o&allTypes != 0
}

func sprintNil(o opts, r reflect.Value) string {
	if _, _, a := withType(o); !a {
		return "nil"
	}

	return fmt.Sprintf("%s(nil)", sprintType(r.Type()))
}

func sprintPrimitive(o opts, r reflect.Value, v interface{}, suppressType ...string) string {
	s := fmt.Sprint(v)
	if s[0] == '(' && s[len(s)-1] == ')' {
		s = s[1 : len(s)-1]
	}

	_, w, a := withType(o)
	if !w {
		return s
	}

	t := sprintType(r.Type())
	if !a {
		for _, suppress := range suppressType {
			if t == suppress {
				return s
			}
		}
	}

	return fmt.Sprintf("%s(%s)", t, s)
}

func sprintItems(o opts, prefix string, r reflect.Value) string {
	o, w, _ := withType(o)
	itemOpts := o | skipTypes
	s := make([]string, r.Len())
	for i := 0; i < r.Len(); i++ {
		s[i] = sprint(itemOpts, r.Index(i))
	}

	if !w {
		return fmt.Sprintf("%s{%s}", prefix, strings.Join(s, ", "))
	}

	return fmt.Sprintf("%s{%s}", sprintType(r.Type()), strings.Join(s, ", "))
}

func sprintHidden(o opts, r reflect.Value, hidden string) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	_, w, _ := withType(o)
	if !w {
		return hidden
	}

	return sprintType(r.Type())
}

func sprintArray(o opts, r reflect.Value) string {
	return sprintItems(o, fmt.Sprintf("[%d]", r.Len()), r)
}

func sprintChan(o opts, r reflect.Value) string {
	return sprintHidden(o, r, "chan")
}

func sprintFunc(o opts, r reflect.Value) string {
	return sprintHidden(o, r, "func()")
}

func sprinterface(o opts, r reflect.Value) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	o, w, _ := withType(o)
	if !w {
		return sprint(o, r.Elem())
	}

	return fmt.Sprintf("%s(%s)", sprintType(r.Type()), sprint(o, r.Elem()))
}

func sprintMap(o opts, r reflect.Value) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	o, w, _ := withType(o)
	itemOpts := o | skipTypes
	var items []string
	for _, key := range r.MapKeys() {
		items = append(
			items,
			fmt.Sprintf(
				"%s: %s",
				sprint(itemOpts, key),
				sprint(itemOpts, r.MapIndex(key)),
			),
		)
	}

	sort.Strings(items)
	sitems := strings.Join(items, ", ")
	if !w {
		return fmt.Sprintf("map{%s}", sitems)
	}

	return fmt.Sprintf("%s{%s}", sprintType(r.Type()), sitems)
}

func sprintPointer(o opts, r reflect.Value) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	s := sprint(o, r.Elem())
	if _, w, _ := withType(o); !w {
		return s
	}

	return fmt.Sprintf("*%s", s)
}

func sprintList(o opts, r reflect.Value) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	return sprintItems(o, "[]", r)
}

func sprintString(o opts, r reflect.Value) string {
	b := []byte(r.String())
	e := make([]byte, 0, len(b))
	for _, c := range b {
		switch c {
		case '\\':
			e = append(e, '\\', '\\')
		case '"':
			e = append(e, '\\', '"')
		case '\b':
			e = append(e, '\\', 'b')
		case '\f':
			e = append(e, '\\', 'f')
		case '\n':
			e = append(e, '\\', 'n')
		case '\r':
			e = append(e, '\\', 'r')
		case '\t':
			e = append(e, '\\', 't')
		case '\v':
			e = append(e, '\\', 'v')
		default:
			e = append(e, c)
		}
	}

	s := fmt.Sprintf("\"%s\"", string(e))
	_, w, a := withType(o)
	if !w {
		return s
	}

	t := sprintType(r.Type())
	if !a && t == "string" {
		return s
	}

	return fmt.Sprintf("%s(%s)", t, s)
}

func sprintStruct(o opts, r reflect.Value) string {
	o, w, _ := withType(o)
	fieldOpts := o | skipTypes
	rt := r.Type()
	f := make([]string, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		name := rt.Field(i).Name
		f[i] = fmt.Sprintf(
			"%s: %s",
			name,
			sprint(fieldOpts, r.FieldByName(name)),
		)
	}

	fs := strings.Join(f, ", ")
	if !w {
		return fmt.Sprintf("{%s}", fs)
	}

	return fmt.Sprintf("%s{%s}", sprintType(rt), fs)
}

func sprintUnsafePointer(o opts, r reflect.Value) string {
	if r.IsNil() {
		return sprintNil(o, r)
	}

	return "pointer"
}

func sprint(o opts, r reflect.Value) string {
	switch r.Kind() {
	case reflect.Bool:
		return sprintPrimitive(o, r, r.Bool(), "bool")
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return sprintPrimitive(o, r, r.Int(), "int")
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return sprintPrimitive(o, r, r.Uint())
	case reflect.Float32, reflect.Float64:
		return sprintPrimitive(o, r, r.Float())
	case reflect.Complex64, reflect.Complex128:
		return sprintPrimitive(o, r, r.Complex())
	case reflect.Array:
		return sprintArray(o, r)
	case reflect.Chan:
		return sprintChan(o, r)
	case reflect.Func:
		return sprintFunc(o, r)
	case reflect.Interface:
		return sprinterface(o, r)
	case reflect.Map:
		return sprintMap(o, r)
	case reflect.Ptr:
		return sprintPointer(o, r)
	case reflect.Slice:
		return sprintList(o, r)
	case reflect.String:
		return sprintString(o, r)
	case reflect.Struct:
		return sprintStruct(o, r)
	case reflect.UnsafePointer:
		return sprintUnsafePointer(o, r)
	default:
		return "<invalid>"
	}
}
