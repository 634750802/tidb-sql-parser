package main

import (
	"encoding/json"
	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"syscall/js"
)

func parse(sql string) string {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(stmtNodes)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func main() {
	// https://zmis.me/user/zmisgod/post/1607
	c := make(chan struct{}, 0)

	jsonParse := js.Global().Get("JSON")

	js.Global().Set("__tidbSqlParse", js.FuncOf(func(this js.Value, args []js.Value) any {
		return jsonParse.Call("parse", js.ValueOf(parse(args[0].String())))
	}))

	<-c
}
