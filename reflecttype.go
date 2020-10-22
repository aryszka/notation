package notation

import (
	"reflect"
)

func reflectFuncBaseType(t reflect.Type) node {
	args := func(num func() int, typ func(int) reflect.Type) []node {
		var t []node
		for i := 0; i < num(); i++ {
			t = append(t, reflectType(typ(i)))
		}

		return t
	}

	in := args(t.NumIn, t.In)
	out := args(t.NumOut, t.Out)

	n := nodeOf("(")
	if len(in) == 1 {
		n.parts = append(n.parts, in[0])
	} else if len(in) > 1 {
		n.parts = append(n.parts, wrapper{sep: ", ", items: in})
	}

	n.parts = append(n.parts, ")")
	if len(out) == 1 {
		n.parts = append(n.parts, " ", out[0])
	} else if len(out) > 1 {
		n.parts = append(n.parts, " (", wrapper{sep: ", ", items: out}, ")")
	}

	return n
}

func reflectArrayType(t reflect.Type) node {
	return nodeOf("[", t.Len(), "]", reflectType(t.Elem()))
}

func reflectChanType(t reflect.Type) node {
	var prefix string
	switch t.ChanDir() {
	case reflect.RecvDir:
		prefix = "<-chan "
	case reflect.SendDir:
		prefix = "chan<- "
	default:
		prefix = "chan "
	}

	return nodeOf(prefix, reflectType(t.Elem()))
}

func reflectFuncType(t reflect.Type) node {
	return nodeOf("func", reflectFuncBaseType(t))
}

func reflectInterfaceType(t reflect.Type) node {
	wr := wrapper{sep: "; "}
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		wr.items = append(
			wr.items,
			nodeOf(
				method.Name,
				reflectFuncBaseType(method.Type),
			),
		)
	}

	return nodeOf("interface{", wr, "}")
}

func reflectMapType(t reflect.Type) node {
	return nodeOf("map[", reflectType(t.Key()), "]", reflectType(t.Elem()))
}

func reflectPointerType(t reflect.Type) node {
	return nodeOf("*", reflectType(t.Elem()))
}

func reflectListType(t reflect.Type) node {
	return nodeOf("[]", reflectType(t.Elem()))
}

func reflectStructType(t reflect.Type) node {
	wr := wrapper{sep: "; "}
	for i := 0; i < t.NumField(); i++ {
		fi := t.Field(i)
		wr.items = append(
			wr.items,
			nodeOf(
				fi.Name,
				" ",
				reflectType(fi.Type),
			),
		)
	}

	return nodeOf("struct{", wr, "}")
}

func reflectType(t reflect.Type) node {
	if t.Name() != "" {
		return nodeOf(t.Name())
	}

	switch t.Kind() {
	case reflect.Array:
		return reflectArrayType(t)
	case reflect.Chan:
		return reflectChanType(t)
	case reflect.Func:
		return reflectFuncType(t)
	case reflect.Interface:
		return reflectInterfaceType(t)
	case reflect.Map:
		return reflectMapType(t)
	case reflect.Ptr:
		return reflectPointerType(t)
	case reflect.Slice:
		return reflectListType(t)
	case reflect.Struct:
		return reflectStructType(t)
	default:
		return nodeOf("<invalid>")
	}
}
