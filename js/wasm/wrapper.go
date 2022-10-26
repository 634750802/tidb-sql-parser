package wasm

import (
	"reflect"
	"syscall/js"
	"tidb-sql-parser/utils"
)

var errSingleResultMessage = utils.RuntimeError{"Calling go function only support on return value. (or second is error)"}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func (r *Runtime) wrapFunc(f any) js.Func {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic("wrapFunc requires a func")
	}

	funcInstance := reflect.ValueOf(f)

	return js.FuncOf(func(this js.Value, args []js.Value) any {
		goArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			goArgs[i] = r.jsValueToValue(arg, t.In(i))
		}

		returned := funcInstance.Call(goArgs)

		val, err := r.handleResult(t, returned)

		if err != nil {
			panic(err)
		}

		return r.valueToJsValue(val)
	})
}

func (r *Runtime) jsValueToValue(arg js.Value, in reflect.Type) reflect.Value {
	switch in.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		return reflect.ValueOf(arg.Int()).Convert(in)
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(arg.Float()).Convert(in)
	case reflect.Bool:
		return reflect.ValueOf(arg.Bool()).Convert(in)
	case reflect.String:
		return reflect.ValueOf(arg.String()).Convert(in)
	case reflect.Pointer, reflect.Interface:
		if arg.Type() == js.TypeNumber {
			i, _ := r.pool.get(ObjectHandle(arg.Int()))
			return reflect.ValueOf(i).Convert(in)
		}
		if in.Elem().Kind() == reflect.Struct {
			s := reflect.New(in.Elem())
			utils.ForEachTypeFields(in.Elem(), func(field reflect.StructField, index []int) {
				s.Elem().FieldByIndex(index).Set(r.jsValueToValue(arg.Get(field.Name), field.Type))
			})
			return s
		}
		break
	case reflect.Struct:
		if in.AssignableTo(reflect.TypeOf(js.Value{})) {
			return reflect.ValueOf(arg).Convert(in)
		}
	default:
	}
	panic("cdo not support arg type " + in.String())
}

func (r *Runtime) valueToJsValue(value reflect.Value) js.Value {
	if !value.IsValid() {
		return js.Undefined()
	}
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return js.ValueOf(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return js.ValueOf(value.Uint())
	case reflect.Float32, reflect.Float64:
		return js.ValueOf(value.Float())
	case reflect.Bool:
		return js.ValueOf(value.Bool())
	case reflect.String:
		return js.ValueOf(value.String())
	case reflect.Interface, reflect.Pointer:
		otp, ok := r.pool.getDefinedType(value.Type())
		if !ok {
			otp = r.pool.tid
			r.pool.defineType(otp, value.Type())
		}
		return r.createRef(r.pool.new(value.Interface()))
	case reflect.Slice:
		jsArr := make([]any, value.Len())
		for i := 0; i < value.Len(); i++ {
			jsArr[i] = r.valueToJsValue(value.Index(i))
		}
		return js.ValueOf(jsArr)
	case reflect.Struct:
		if value.Type().AssignableTo(reflect.TypeOf(js.Value{})) {
			return js.ValueOf(value.Interface())
		}
		res := map[string]any{}

		utils.ForEachTypeFields(value.Type(), func(field reflect.StructField, index []int) {
			res[field.Name] = r.valueToJsValue(value.FieldByIndex(index))
		})
		return js.ValueOf(res)
	default:
		panic("do not support type " + value.Type().String())
	}
}

func (r *Runtime) handleResult(mt reflect.Type, result []reflect.Value) (reflect.Value, error) {
	no := mt.NumOut()
	if no > 2 {
		return reflect.ValueOf(nil), &errSingleResultMessage
	}
	if no == 2 {
		if !mt.Out(1).Implements(errorInterface) {
			return reflect.ValueOf(nil), &errSingleResultMessage
		}
		if result[1].IsNil() {
			return result[0], nil
		} else {
			return reflect.ValueOf(nil), result[1].Interface().(error)
		}
	}
	if no == 1 {
		return result[0], nil
	} else {
		return reflect.ValueOf(nil), nil
	}
}
