package analyze

import (
	"encoding/json"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/types"
)

type JsonType byte

type Tp struct {
	Type     types.EvalType
	Nullable bool
}

const (
	ETInt       = "ETInt"
	ETReal      = "ETReal"
	ETDecimal   = "ETDecimal"
	ETString    = "ETString"
	ETDatetime  = "ETDatetime"
	ETTimestamp = "ETTimestamp"
	ETDuration  = "ETDuration"
	ETJson      = "ETJson"
)

var EvalTypes = []string{ETInt, ETReal, ETDecimal, ETString, ETDatetime, ETTimestamp, ETDuration, ETJson}

func (t *Tp) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Nullable bool
	}{
		EvalTypes[int(t.Type)],
		t.Nullable,
	})
}

type Column struct {
	Tp
	Name string
}

func NewColumn(name string, tp types.EvalType, nullable bool) *Column {
	return &Column{Tp: Tp{tp, nullable}, Name: name}
}

func (c Column) as(name string) *Column {
	if name == "" {
		return &c
	} else {
		return &Column{
			Tp: Tp{
				Type:     c.Type,
				Nullable: c.Nullable,
			},
			Name: name,
		}
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
