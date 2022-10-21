package analyze

import (
	"github.com/pingcap/tidb/parser/types"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)
import _ "github.com/pingcap/tidb/parser/test_driver"

func TestSimpleSQL(t *testing.T) {
	p := NewParser()

	p.DefineFunc("CURRENT_TIMESTAMP", &Tp{types.ETTimestamp, false})
	p.AddDdl("create table test(id integer primary key, name varchar(255) not null);")

	id := NewColumn("id", types.ETInt, true)
	name := NewColumn("name", types.ETString, false)

	AssertColumnEquals(t, p.ctx.getTable("test").GetColumn("id"), id)
	AssertColumnEquals(t, p.ctx.getTable("test").GetColumn("name"), name)

	columns := p.Parse("select id, t.name, 1 as d, 's' as s, CURRENT_TIMESTAMP as t, 1 + 1.2 as n, 1 + 1 from test as t")

	AssertColumnEquals(t, columns[0], id)
	AssertColumnEquals(t, columns[1], name)
	AssertColumnEquals(t, columns[2], NewColumn("d", types.ETInt, false))
	AssertColumnEquals(t, columns[3], NewColumn("s", types.ETString, false))
	AssertColumnEquals(t, columns[4], NewColumn("t", types.ETTimestamp, false))
	AssertColumnEquals(t, columns[5], NewColumn("n", types.ETInt, false))
	AssertColumnEquals(t, columns[6], NewColumn("1 + 1", types.ETInt, false))
}

func readTestResource(t *testing.T, file string) string {
	_, filename, _, _ := runtime.Caller(0)

	bytes, err := os.ReadFile(filepath.Dir(filename) + "/" + file)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func TestTrendingRepoSQL(t *testing.T) {
	p := NewParser()

	p.AddDdl(readTestResource(t, "test_schema.sql"))
	p.DefineTransparentFunc("IFNULL")
	p.DefineTransparentFunc("SUM")
	p.DefineTransparentFunc("ABS")
	p.DefineTransparentFunc("GREATEST")
	p.DefineTransparentFunc("LEAST")
	p.DefineFunc("DATE_SUB", &Tp{types.ETDatetime, false})
	p.DefineFunc("COUNT", &Tp{types.ETInt, false})
	p.DefineFunc("TIMESTAMPDIFF", &Tp{types.ETReal, false})
	cols := p.Parse(readTestResource(t, "test_query.sql"))
	for _, col := range cols {
		println(col.Name, col.Type, col.Nullable)
	}
}
