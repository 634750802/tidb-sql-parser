package main

import (
	_ "github.com/pingcap/tidb/parser/test_driver"
	"reflect"
	"tidb-sql-parser/analyze"
	"tidb-sql-parser/js/wasm"
)

type Globals struct {
}

func (*Globals) NewParser() *analyze.Parser {
	return analyze.NewParser()
}

func main() {
	// https://zmis.me/user/zmisgod/post/1607
	c := make(chan struct{}, 0)
	r := wasm.NewRuntime()
	r.DefineType(0, reflect.TypeOf(Globals{}))
	<-c
}
