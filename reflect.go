package notation

import (
	"bytes"
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

func reflectItems(o opts, p *pending, prefix string, r reflect.Value) node {
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
				reflectValue(itemOpts, p, r.Index(i)),
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

func reflectArray(o opts, p *pending, r reflect.Value) node {
	return reflectItems(o, p, fmt.Sprintf("[%d]", r.Len()), r)
}

func reflectChan(o opts, r reflect.Value) node {
	return reflectHidden(o, "chan", r)
}

func reflectFunc(o opts, r reflect.Value) node {
	return reflectHidden(o, "func()", r)
}

func reflectInterface(o opts, p *pending, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, false, r)
	}

	e := reflectValue(o&^skipTypes, p, r.Elem())
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

func reflectMap(o opts, p *pending, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	var (
		nkeys []node
		skeys []string
	)

	// TODO: simplify this when no sorting is required
	items := wrapper{sep: ", ", suffix: ","}
	itemOpts := o | skipTypes
	keys := r.MapKeys()
	sv := make(map[string]reflect.Value)
	sn := make(map[string]node)
	for _, key := range keys {
		var b bytes.Buffer
		nk := reflectValue(itemOpts, p, key)
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
				reflectValue(
					itemOpts,
					p,
					r.MapIndex(sv[skey]),
				),
			),
		)
	}

	if _, t, _ := withType(o); !t {
		return nodeOf("map{", items, "}")
	}

	return nodeOf(reflectType(r.Type()), "{", items, "}")
}

func reflectPointer(o opts, p *pending, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	e := reflectValue(o, p, r.Elem())
	if _, t, _ := withType(o); !t {
		return e
	}

	return nodeOf("*", e)
}

func reflectList(o opts, p *pending, r reflect.Value) node {
	if r.IsNil() {
		return reflectNil(o, true, r)
	}

	return reflectItems(o, p, "[]", r)
}

func reflectString(o opts, r reflect.Value) node {
	sv := r.String()
	b := []byte(sv)
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

	s := str{val: fmt.Sprintf("\"%s\"", string(e))}
	if !strings.Contains(sv, "`") && strings.Contains(sv, "\n") {
		s.raw = fmt.Sprintf("`%s`", sv)
	}

	n := nodeOf(s)
	_, t, a := withType(o)
	if !t {
		return n
	}

	tn := reflectType(r.Type())
	if !a && tn.parts[0] == "string" {
		return n
	}

	return nodeOf(tn, "(", wrapper{items: []node{n}}, ")")
}

func reflectStruct(o opts, p *pending, r reflect.Value) node {
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
					p,
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

func checkPending(p *pending, r reflect.Value) (applyRef func(node) node, ref node, isPending bool) {
	applyRef = func(n node) node { return n }
	switch r.Kind() {
	case reflect.Slice, reflect.Map:
	case reflect.Ptr:
		if r.IsNil() {
			return
		}
	default:
		return
	}

	var nr nodeRef
	key := r.Pointer()
	nr, isPending = p.values[key]
	if isPending {
		nr.refCount++
		p.values[key] = nr
		ref = nodeOf("r", nr.id)
		return
	}

	nr = nodeRef{id: p.idCounter}
	p.idCounter++
	p.values[key] = nr
	applyRef = func(n node) node {
		nr = p.values[key]
		if nr.refCount > 0 {
			n.parts = append(
				[]interface{}{"r", nr.id, "="},
				n.parts...,
			)
		}

		delete(p.values, key)
		return n
	}

	return
}

func reflectValue(o opts, p *pending, r reflect.Value) node {
	applyRef, ref, isPending := checkPending(p, r)
	if isPending {
		return ref
	}

	var n node
	switch r.Kind() {
	case reflect.Bool:
		n = reflectPrimitive(o, r, r.Bool(), "bool")
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		n = reflectPrimitive(o, r, r.Int(), "int")
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		n = reflectPrimitive(o, r, r.Uint())
	case reflect.Float32, reflect.Float64:
		n = reflectPrimitive(o, r, r.Float())
	case reflect.Complex64, reflect.Complex128:
		n = reflectPrimitive(o, r, r.Complex())
	case reflect.Array:
		n = reflectArray(o, p, r)
	case reflect.Chan:
		n = reflectChan(o, r)
	case reflect.Func:
		n = reflectFunc(o, r)
	case reflect.Interface:
		n = reflectInterface(o, p, r)
	case reflect.Map:
		n = reflectMap(o, p, r)
	case reflect.Ptr:
		n = reflectPointer(o, p, r)
	case reflect.Slice:
		n = reflectList(o, p, r)
	case reflect.String:
		n = reflectString(o, r)
	case reflect.UnsafePointer:
		n = reflectUnsafePointer(o, r)
	default:
		n = reflectStruct(o, p, r)
	}

	return applyRef(n)
}
