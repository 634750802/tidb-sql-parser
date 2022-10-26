package wasm

import (
	"reflect"
	"syscall/js"
	"tidb-sql-parser/utils"
)

type GeneralCmd uint32
type ObjectHandle uint32
type ObjectId uint32
type ObjectType uint32

const (
	GeneralCmdNew GeneralCmd = iota
	GeneralCmdDelete
	GeneralCmdReflectGetProperty
	GeneralCmdReflectGetMethod
	GeneralCmdReflectSetProperty
	GeneralCmdReflectCallMethod
	GeneralCmdReflectGetPropertyNames
	GeneralCmdDescribeType
	GeneralCmdExit = 0xffffffff
)

const (
	typeNode   ObjectType = 0x01000000
	typeParser            = 0x02000000
)

const (
	codeOk                 uint32 = 0x00000000
	codeBadType            uint32 = 0x01000000
	codeBadHandle          uint32 = 0x02000000
	codeBadField           uint32 = 0x03000000
	codeBadResult          uint32 = 0x04000000
	codeBadMethodSignature uint32 = 0x05000000
	codeBadCmd             uint32 = 0xff000000
)

const typeMask uint32 = 0xff000000
const idMask uint32 = 0x00ffffff

func (r *Runtime) executeCmd(cmd GeneralCmd) uint32 {
	switch cmd {
	case GeneralCmdNew:
		return r.cmdNew()
	case GeneralCmdDelete:
		return r.cmdDelete()
	case GeneralCmdReflectGetProperty:
		return r.cmdReflectGetProperty()
	case GeneralCmdReflectGetMethod:
		return r.cmdReflectGetMethod()
	case GeneralCmdReflectSetProperty:
		return r.cmdReflectSetProperty()
	case GeneralCmdReflectCallMethod:
		return r.cmdReflectCallMethod()
	case GeneralCmdReflectGetPropertyNames:
		return r.cmdReflectGetPropertyNames()
	case GeneralCmdDescribeType:
		return r.cmdDescribeType()
	}
	return 0xffffffff
}

func getIntSlice(v js.Value) []int {
	l := v.Length()
	s := make([]int, l)
	for i := 0; i < l; i++ {
		s[i] = v.Index(i).Int()
	}
	return s
}

func (r *Runtime) cmdNew() uint32 {
	handle := r.getCtxHandle()
	otp := ObjectType(uint32(handle) & typeMask)
	if tp := r.pool.getType(otp); tp != nil {
		v := reflect.New(tp)
		return uint32(r.pool.new(v.Interface())) & idMask
	}
	r.setCtxErrMsg("Unknown type")
	return codeBadType
}

func (r *Runtime) cmdDelete() uint32 {
	handle := r.getCtxHandle()
	return uint32(r.pool.delete(handle))
}

func (r *Runtime) cmdReflectGetProperty() uint32 {
	handle := r.getCtxHandle()
	pn := getIntSlice(r.getCtxArg(0))

	obj, _ := r.pool.get(handle)

	if obj == nil {
		r.setCtxErrMsg("Bad object handle")
		return codeBadHandle
	}

	r.setCtxRet(r.valueToJsValue(reflect.ValueOf(obj).Elem().FieldByIndex(pn)))
	return codeOk
}

func (r *Runtime) cmdReflectGetMethod() uint32 {
	handle := r.getCtxHandle()
	pn := r.getCtxArg(0).Int()

	obj, tp := r.pool.get(handle)

	m := tp.Method(pn)

	r.setCtxRet(js.FuncOf(func(this js.Value, args []js.Value) any {
		fArgs := make([]reflect.Value, m.Type.NumIn())
		fArgs[0] = reflect.ValueOf(obj)
		for i := 1; i < m.Type.NumIn(); i++ {
			fArgs[i] = r.jsValueToValue(args[i-1], m.Type.In(i))
		}
		ret := m.Func.Call(fArgs)

		v, err := r.handleResult(m.Type, ret)
		if err != nil {
			panic(err)
		}

		return r.valueToJsValue(v)
	}))
	return codeOk
}

func (r *Runtime) cmdReflectSetProperty() uint32 {
	handle := r.getCtxHandle()
	pn := getIntSlice(r.getCtxArg(0))
	rpv := r.getCtxArg(1)

	obj, tp := r.pool.get(handle)

	if obj == nil {
		r.setCtxErrMsg("Bad object handle")
		return codeBadHandle
	}

	f := tp.FieldByIndex(pn)

	pv := r.jsValueToValue(rpv, f.Type)

	reflect.ValueOf(obj).Elem().FieldByIndex(pn).Set(pv)
	return codeOk
}

func (r *Runtime) cmdReflectCallMethod() uint32 {
	handle := r.getCtxHandle()
	fn := r.getCtxArg(0).Int()
	args := r.getCtxArg(1)

	obj, tp := r.pool.get(handle)

	if obj == nil {
		r.setCtxErrMsg("Bad object handle")
		return codeBadHandle
	}

	m := reflect.PointerTo(tp).Method(fn)
	mt := m.Type

	n := mt.NumIn() - 1
	in := make([]reflect.Value, n)
	i := 0
	for i = 0; i < n; i++ {
		in[i] = r.jsValueToValue(args.Index(i), mt.In(i+1))
	}

	ret := reflect.ValueOf(obj).Method(m.Index).Call(in)

	v, err := r.handleResult(mt, ret)
	if err != nil {
		r.setCtxErrMsg(err.Error())
		return codeBadCmd
	}
	r.setCtxRet(r.valueToJsValue(v))
	return codeOk
}

func (r *Runtime) cmdReflectGetPropertyNames() uint32 {
	handle := r.getCtxHandle()

	obj, tp := r.pool.get(handle)
	if obj == nil {
		r.setCtxErrMsg("Bad object handle")
		return codeBadHandle
	}

	ri := 0
	el := tp.Elem()
	names := make([]any, el.NumField()+tp.NumMethod())
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		if f.IsExported() {
			names[ri] = f.Name
			ri++
		}
	}

	ri = 0
	off := el.NumField()
	for i := 0; i < tp.NumMethod(); i++ {
		m := tp.Method(i)
		if m.IsExported() {
			names[off+ri] = m.Name
			ri++
		}
	}

	r.setCtxRet(names)
	return codeOk
}

func (r *Runtime) cmdDescribeType() uint32 {
	handle := r.getCtxHandle()

	if tp := r.pool.getType(ObjectType(uint32(handle) & typeMask)); tp != nil {
		fields := js.ValueOf(map[string]any{})
		methods := js.ValueOf(map[string]any{})
		info := js.ValueOf(map[string]any{
			"id":      tp.String(),
			"size":    tp.Size(),
			"fields":  fields,
			"methods": methods,
		})

		ri := 0
		utils.ForEachTypeFields(tp, func(field reflect.StructField, index []int) {
			fields.Set(field.Name, utils.EraseType(index))
			ri++
		})

		ri = 0
		ptr := reflect.PointerTo(tp)
		utils.ForEachTypeMethods(ptr, func(method reflect.Method, index int) {
			methods.Set(method.Name, index)
			ri++
		})

		r.setCtxRet(info)
		return codeOk
	}
	return codeBadType
}
