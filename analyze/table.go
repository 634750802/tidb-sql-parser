package analyze

import (
	"encoding/json"
	"strings"
)

const UNNAMED = "_unnamed_"

type Table interface {
	Name() string
	Columns() []*Column

	GetColumn(name string) *Column
}

type TableDefine struct {
	columnsArr []*Column
	columns    map[string]*Column
	name       string
}

func (t *TableDefine) Columns() []*Column {
	return t.columnsArr
}

func (t *TableDefine) Name() string {
	return t.name
}

func NewTableDefine(name string) *TableDefine {
	return &TableDefine{
		columns: make(map[string]*Column),
		name:    name,
	}
}

func (t *TableDefine) AddColumn(column *Column) {
	t.columns[column.Name] = column
	t.columnsArr = append(t.columnsArr, column)
}

func (t *TableDefine) GetColumn(name string) *Column {
	return t.columns[strings.ToLower(name)]
}

func (t *TableDefine) Merge(tables ...Table) {
	for _, table := range tables {
		for _, column := range table.Columns() {
			if t.GetColumn(column.Name) == nil {
				t.AddColumn(column)
			}
		}
	}
}

func (t *TableDefine) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string
		Columns []*Column
	}{
		Name:    t.name,
		Columns: t.columnsArr,
	})
}

func MergeTables(as string, tables ...Table) Table {
	res := NewTableDefine(as)

	for _, table := range tables {
		for _, column := range table.Columns() {
			if res.GetColumn(column.Name) == nil {
				res.AddColumn(column)
			}
		}
	}

	return res
}

type TableRef struct {
	as    string
	table Table
}

func NewTableRef(as string, table Table) Table {
	if as == UNNAMED {
		return table
	} else {
		return &TableRef{as: as, table: table}
	}
}

func (t *TableRef) Name() string {
	return t.as
}

func (t *TableRef) Columns() []*Column {
	return t.table.Columns()
}

func (t *TableRef) GetColumn(name string) *Column {
	return t.table.GetColumn(name)
}

func (t *TableRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string
		As      string
		Columns []*Column
	}{
		Name:    t.Name(),
		As:      t.Name(),
		Columns: t.Columns(),
	})
}
