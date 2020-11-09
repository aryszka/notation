package notation

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
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

func reflectPrimitive(o opts, r reflect.Value, v interface{}, suppressType ...string) node {
	s := fmt.Sprint(v)
	if s[0] == '(' && s[len(s)-1] == ')' {
		s = s[1 : len(s)-1]
	}

	_, t, a := withType(o)
	if !t {
		return nodeOf(s)
	}

	tn := reflectType(r.Type())
	if a {
		return nodeOf(tn, "(", s, ")")
	}

	for _, suppress := range suppressType {
		if tn.parts[0] != suppress {
			continue
		}

		return nodeOf(s)
	}

	return nodeOf(tn, "(", s, ")")
}

func reflectNil(o opts, groupUnnamedType bool, r reflect.Value) node {
	if _, _, a := withType(o); !a {
		return nodeOf("nil")
	}

	rt := r.Type()
	if groupUnnamedType && rt.Name() == "" {
		return nodeOf("(", reflectType(rt), ")(nil)")
	}

	return nodeOf(reflectType(rt), "(nil)")
}

func reflectItems(o opts, prefix string, r reflect.Value) node {
	typ := r.Type()
	var items wrapper
	if typ.Elem().Name() == "uint8" {
		items = wrapper{sep: " ", mode: line}
		for i := 0; i < r.Len(); i++ {
			items.items = append(
				items.items,
				nodeOf(fmt.Sprintf("%02x", r.Index(i).Uint())),
			)
		}
	} else {
		items = wrapper{sep: ", ", suffix: ","}
		itemOpts := o | skipTypes
		for i := 0; i < r.Len(); i++ {
			items.items = append(
				items.items,
				reflectValue(itemOpts, r.Index(i)),
			)
		}
	}

	if _, t, _ := withType(o); !t {
		return nodeOf(prefix, "{", items, "}")
	}

	return nodeOf(reflectType(typ), "{", items, "}")
}

func reflectHidden(o opts, hidden string, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	if _, t, _ := withType(o); !t {
		return nodeOf(hidden)
	}

	return reflectType(r.Type())
}

func reflectArray(o opts, r reflect.Value) node {
	return reflectItems(o, fmt.Sprintf("[%d]", r.Len()), r)
}

func reflectChan(o opts, r reflect.Value) node {
	return reflectHidden(o, "chan", r)
}

func reflectFunc(o opts, r reflect.Value) node {
	return reflectHidden(o, "func()", r)
}

func reflectInterface(o opts, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, false, r)
	}

	e := reflectValue(o, r.Elem())
	if _, t, _ := withType(o); !t {
		return e
	}

	return nodeOf(
		reflectType(r.Type()),
		"(",
		wrapper{items: []node{e}},
		")",
	)
}

func reflectMap(o opts, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	var (
		nkeys []node
		skeys []string
	)

	items := wrapper{sep: ", ", suffix: ","}
	itemOpts := o | skipTypes
	keys := r.MapKeys()
	sv := make(map[string]reflect.Value)
	sn := make(map[string]node)
	for _, key := range keys {
		var b bytes.Buffer
		nk := reflectValue(itemOpts, key)
		nkeys = append(nkeys, nk)
		wr := writer{w: &b}
		fprint(&wr, 0, nk)
		skey := b.String()
		skeys = append(skeys, skey)
		sv[skey] = key
		sn[skey] = nk
	}

	if o&randomMaps == 0 {
		sort.Strings(skeys)
	}

	for _, skey := range skeys {
		items.items = append(
			items.items,
			nodeOf(
				sn[skey],
				": ",
				reflectValue(itemOpts, r.MapIndex(sv[skey])),
			),
		)
	}

	if _, t, _ := withType(o); !t {
		return nodeOf("map{", items, "}")
	}

	return nodeOf(reflectType(r.Type()), "{", items, "}")
}

func reflectPointer(o opts, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	e := reflectValue(o, r.Elem())
	if _, t, _ := withType(o); !t {
		return e
	}

	return nodeOf("*", e)
}

func reflectList(o opts, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	return reflectItems(o, "[]", r)
}

func reflectString(o opts, r reflect.Value) node {
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
	_, t, a := withType(o)
	if !t {
		return nodeOf(s)
	}

	tn := reflectType(r.Type())
	if !a && tn.parts[0] == "string" {
		return nodeOf(s)
	}

	return nodeOf(tn, "(", wrapper{items: []node{nodeOf(s)}}, ")")
}

func reflectStruct(o opts, r reflect.Value) node {
	wr := wrapper{sep: ", ", suffix: ","}

	fieldOpts := o | skipTypes
	rt := r.Type()
	for i := 0; i < r.NumField(); i++ {
		name := rt.Field(i).Name
		wr.items = append(
			wr.items,
			nodeOf(
				name,
				": ",
				reflectValue(
					fieldOpts,
					r.FieldByName(name),
				),
			),
		)
	}

	if _, t, _ := withType(o); !t {
		return nodeOf("{", wr, "}")
	}

	return nodeOf(reflectType(rt), "{", wr, "}")
}

func reflectUnsafePointer(o opts, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, false, r)
	}

	if _, _, a := withType(o); !a {
		return nodeOf("pointer")
	}

	return nodeOf(reflectType(r.Type()), "(pointer)")
}

func reflectValue(o opts, r reflect.Value) node {
	switch r.Kind() {
	case reflect.Bool:
		return reflectPrimitive(o, r, r.Bool(), "bool")
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return reflectPrimitive(o, r, r.Int(), "int")
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return reflectPrimitive(o, r, r.Uint())
	case reflect.Float32, reflect.Float64:
		return reflectPrimitive(o, r, r.Float())
	case reflect.Complex64, reflect.Complex128:
		return reflectPrimitive(o, r, r.Complex())
	case reflect.Array:
		return reflectArray(o, r)
	case reflect.Chan:
		return reflectChan(o, r)
	case reflect.Func:
		return reflectFunc(o, r)
	case reflect.Interface:
		return reflectInterface(o, r)
	case reflect.Map:
		return reflectMap(o, r)
	case reflect.Ptr:
		return reflectPointer(o, r)
	case reflect.Slice:
		return reflectList(o, r)
	case reflect.String:
		return reflectString(o, r)
	case reflect.UnsafePointer:
		return reflectUnsafePointer(o, r)
	default:
		return reflectStruct(o, r)
	}
}
