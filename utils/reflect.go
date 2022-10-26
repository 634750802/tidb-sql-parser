package utils

import "reflect"

func _forEachTypeFields(tp reflect.Type, index []int, handle func(field reflect.StructField, index []int)) {
	n := tp.NumField()
	for i := 0; i < n; i++ {
		f := tp.Field(i)

		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			_forEachTypeFields(f.Type, append(index, i), handle)
		} else {
			handle(f, append(index, i))
		}
	}
}

func ForEachTypeFields(tp reflect.Type, handle func(field reflect.StructField, index []int)) {
	_forEachTypeFields(tp, []int{}, handle)
}

func ForEachTypeMethods(tp reflect.Type, handle func(method reflect.Method, index int)) {
	n := tp.NumMethod()
	for i := 0; i < n; i++ {
		m := tp.Method(i)

		if !m.IsExported() {
			continue
		}

		handle(m, i)
	}
}
