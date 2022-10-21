package analyze

import (
	"github.com/pingcap/tidb/parser/types"
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

	columns := p.Parse("select id, name, 1 as d, 's' as s, CURRENT_TIMESTAMP as t, 1 + 1.2 as n, 1 + 1 from test")

	AssertColumnEquals(t, columns[0], id)
	AssertColumnEquals(t, columns[1], name)
	AssertColumnEquals(t, columns[2], NewColumn("d", types.ETInt, false))
	AssertColumnEquals(t, columns[3], NewColumn("s", types.ETString, false))
	AssertColumnEquals(t, columns[4], NewColumn("t", types.ETTimestamp, false))
	AssertColumnEquals(t, columns[5], NewColumn("n", types.ETInt, false))
	AssertColumnEquals(t, columns[6], NewColumn("1 + 1", types.ETInt, false))
}
