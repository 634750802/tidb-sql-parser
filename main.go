package main

import (
	"encoding/json"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/types"
	"syscall/js"
	"tidb-sql-parser/analyze"
	"tidb-sql-parser/utils"
)

const (
	BridgeCmdParserIsOpen int = iota
	BridgeCmdParserOpen
	BridgeCmdParserClose
	BridgeCmdParserAddDdl
	BridgeCmdParserDefineFunc
	BridgeCmdParserParse
	BridgeCmdParserGetTable
	BridgeCmdParserDefineTransparentFunc
	BridgeCmdParserStaticNormalizeDigest
)

const (
	BridgeCodeOk int = iota
	BridgeCodeBadCmd
	BridgeCodeParserAlreadyOpen
	BridgeCodeParserNotOpen
)

type cmdReturn struct {
	Code int
	Data any
}

type cmdCtx struct {
	p *analyze.Parser
}

var ok = &cmdReturn{BridgeCodeOk, "OK"}
var gid = 0
var parsers = map[int]*analyze.Parser{}

func cmd(id int, p *analyze.Parser, cmd int, args []js.Value) *cmdReturn {
	switch cmd {
	case BridgeCmdParserIsOpen:
		return &cmdReturn{BridgeCodeOk, p != nil}
	case BridgeCmdParserOpen:
		if p != nil {
			return &cmdReturn{BridgeCodeParserAlreadyOpen, "BridgeCodeParserAlreadyOpen"}
		}
		id = gid + 1
		gid = id
		parsers[id] = analyze.NewParser()
		return &cmdReturn{BridgeCodeOk, id}
	case BridgeCmdParserClose:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		parsers[id] = nil
		return ok
	case BridgeCmdParserAddDdl:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		p.AddDdl(args[0].String())
		return ok
	case BridgeCmdParserDefineFunc:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		p.DefineFunc(args[0].String(), &analyze.Tp{
			Type:     types.EvalType(args[1].Int()),
			Nullable: args[2].Bool(),
		})
		return ok
	case BridgeCmdParserParse:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		cols := p.Parse(args[0].String())
		return &cmdReturn{BridgeCodeOk, cols}
	case BridgeCmdParserGetTable:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		return &cmdReturn{BridgeCodeOk, p.GetTable(args[0].String())}
	case BridgeCmdParserDefineTransparentFunc:
		if p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		p.DefineTransparentFunc(args[0].String())
		return ok
	case BridgeCmdParserStaticNormalizeDigest:
		return &cmdReturn{BridgeCodeOk, analyze.NormalizeDigest(args[0].String())}
	default:
		return &cmdReturn{BridgeCodeBadCmd, "BridgeCodeBadCmd"}
	}
}

func main() {
	defer func() {
		js.Global().Delete("__tidbSqlParserExecuteCmd")
		js.Global().Delete("__tidbSqlParserEvalTypes")
	}()

	// https://zmis.me/user/zmisgod/post/1607
	c := make(chan struct{}, 0)

	jsonParse := js.Global().Get("JSON")

	js.Global().Set("__tidbSqlParserExecuteCmd", js.FuncOf(func(this js.Value, args []js.Value) any {
		id := this.Int()
		var res any

		p := parsers[id]
		res = cmd(id, p, args[0].Int(), args[1:])

		j, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		return jsonParse.Call("parse", string(j))
	}))
	js.Global().Set("__tidbSqlParserEvalTypes", js.ValueOf(utils.EraseType(analyze.EvalTypes)))
	<-c
}
