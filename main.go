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
var ctx = &cmdCtx{}

func (c *cmdCtx) cmd(cmd int, args []js.Value) *cmdReturn {
	switch cmd {
	case BridgeCmdParserIsOpen:
		return &cmdReturn{BridgeCodeOk, c.p != nil}
	case BridgeCmdParserOpen:
		if c.p != nil {
			return &cmdReturn{BridgeCodeParserAlreadyOpen, "BridgeCodeParserAlreadyOpen"}
		}
		c.p = analyze.NewParser()
		return ok
	case BridgeCmdParserClose:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		c.p = nil
		return ok
	case BridgeCmdParserAddDdl:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		c.p.AddDdl(args[0].String())
		return ok
	case BridgeCmdParserDefineFunc:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		c.p.DefineFunc(args[0].String(), &analyze.Tp{
			Type:     types.EvalType(args[1].Int()),
			Nullable: args[2].Bool(),
		})
		return ok
	case BridgeCmdParserParse:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		cols := c.p.Parse(args[0].String())
		return &cmdReturn{BridgeCodeOk, cols}
	case BridgeCmdParserGetTable:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		return &cmdReturn{BridgeCodeOk, c.p.GetTable(args[0].String())}
	case BridgeCmdParserDefineTransparentFunc:
		if c.p == nil {
			return &cmdReturn{BridgeCodeParserNotOpen, "BridgeCodeParserNotOpen"}
		}
		c.p.DefineTransparentFunc(args[0].String())
		return ok
	case BridgeCmdParserStaticNormalizeDigest:
		return &cmdReturn{BridgeCodeOk, analyze.NormalizeDigest(args[0].String())}
	default:
		return &cmdReturn{BridgeCodeBadCmd, "BridgeCodeBadCmd"}
	}
}

func main() {
	// https://zmis.me/user/zmisgod/post/1607
	c := make(chan struct{}, 0)

	jsonParse := js.Global().Get("JSON")

	js.Global().Set("__tidbSqlParserExecuteCmd", js.FuncOf(func(this js.Value, args []js.Value) any {
		res := ctx.cmd(args[0].Int(), args[1:])
		j, err := json.Marshal(*res)
		if err != nil {
			panic(err)
		}
		return jsonParse.Call("parse", string(j))
	}))
	js.Global().Set("__tidbSqlParserEvalTypes", js.ValueOf(utils.EraseType(analyze.EvalTypes)))
	<-c
}
