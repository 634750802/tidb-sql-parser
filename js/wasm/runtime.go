package wasm

import (
	"os"
	"reflect"
	"syscall/js"
)

type Runtime struct {
	pool *pool

	gCtx       js.Value
	gCreateRef js.Value
}

func NewRuntime() *Runtime {
	r := &Runtime{pool: newPool()}
	r.init()
	return r
}

func (r *Runtime) init() {
	pfx, pfxExists := os.LookupEnv("WASM_PREFIX")
	if !pfxExists {
		panic("WASM_PREFIX not provided")
	}
	global := js.Global()
	global.Set(pfx+"init", r.wrapFunc(func(ctx js.Value, createRef js.Value) {
		r.setJsCtx(ctx, createRef)
		global.Delete(pfx + "init")
		global.Set(pfx+"cmd", r.wrapFunc(func(c int) uint32 {
			return r.executeCmd(GeneralCmd(c))
		}))
	}))
	global.Set(pfx+"rt", r.wrapFunc(func() any {
		return r
	}))
}

func (r *Runtime) DefineType(id byte, tp reflect.Type) {
	r.pool.defineType(ObjectType(uint32(id)<<24), tp)
}
