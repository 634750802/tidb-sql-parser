package analyze

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/types"
)

type JsonType byte

type Tp struct {
	Type     types.EvalType
	Nullable bool
}

var EvalTypes = []string{"ETInt", "ETReal", "ETDecimal", "ETString", "ETDatetime", "ETTimestamp", "ETDuration", "ETJson"}

type Column struct {
	Tp
	Name string
	As   string
}

func NewColumn(name string, tp types.EvalType, nullable bool) *Column {
	return &Column{Tp: Tp{tp, nullable}, Name: name}
}

func (c *Column) as(name string) *Column {
	rn := c.Name
	if rn == UNNAMED {
		rn = name
	}
	return &Column{
		Tp:   c.Tp,
		Name: rn,
		As:   name,
	}
}

func parseColumn(def *ast.ColumnDef) *Column {
	nullable := true

	for _, option := range def.Options {
		switch option.Tp {
		case ast.ColumnOptionNotNull:
			nullable = false
		}
	}

	return &Column{
		Tp: Tp{
			Type:     def.Tp.EvalType(),
			Nullable: nullable,
		},
		Name: def.Name.Name.L,
	}
}
