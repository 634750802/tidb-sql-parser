package wasm

import (
	"syscall/js"
	"unsafe"
)

func (r *Runtime) setJsCtx(ctx js.Value, createRef js.Value) {
	r.gCtx = ctx
	r.gCreateRef = createRef
}

func itu(d int) uint {
	ptr := unsafe.Pointer(&d)
	return *((*uint)(ptr))
}

func uti(d uint) int {
	ptr := unsafe.Pointer(&d)
	return *((*int)(ptr))
}

func (r *Runtime) getCtxHandle() ObjectHandle {
	return ObjectHandle(itu(r.gCtx.Index(1).Int()))
}

func (r *Runtime) getCtxArg(i int) js.Value {
	return r.gCtx.Index(2).Index(i)
}

func (r *Runtime) setCtxRet(value any) {
	r.gCtx.SetIndex(3, value)
}

func (r *Runtime) setCtxErrMsg(s string) {
	r.gCtx.SetIndex(4, js.ValueOf(s))
}

func (r *Runtime) createRef(handle ObjectHandle) js.Value {
	return r.gCreateRef.Invoke(js.ValueOf(uint32(handle)))
}
