package analyze

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/types"
	"reflect"
	"strings"
	"tidb-sql-parser/utils"
)

type Parser struct {
	ctx    *ParseContext
	parser *parser.Parser

	FuncTypes map[string]*Tp
}

func NewParser() *Parser {
	p := &Parser{
		ctx:       newParseContext(nil),
		parser:    parser.New(),
		FuncTypes: map[string]*Tp{},
	}
	p.ctx.parser = p
	return p
}

func (p *Parser) DefineFunc(name string, tp *Tp) {
	p.FuncTypes[strings.ToLower(name)] = tp
}

func (p *Parser) FuncType(name string) *Tp {
	return p.FuncTypes[strings.ToLower(name)]
}

func (p *Parser) AddDdl(sql string) {
	stmts, _, err := p.parser.ParseSQL(sql)
	if err != nil {
		panic(err)
	}
	for _, stmt := range stmts {
		if c, ok := stmt.(*ast.CreateTableStmt); ok {
			p.ctx.addTables(parseCreateTableStmt(c))
		} else {
			panic(unimplementedNodeType(stmt))
		}
	}
}

func (p *Parser) Parse(sql string) []*Column {
	stmt, err := p.parser.ParseOneStmt(sql, "", "")
	if err != nil {
		panic(err)
	}
	if s, ok := stmt.(*ast.SelectStmt); ok {
		return p.ctx.parseSelectStmt(UNNAMED, s).Columns()
	} else {
		panic(unimplementedNodeType(stmt))
	}
}

type ParseContext struct {
	parser *Parser
	parent *ParseContext
	tables map[string]Table

	warn []string
}

func newParseContext(parent *ParseContext) *ParseContext {
	if parent == nil {
		return &ParseContext{tables: map[string]Table{}, warn: make([]string, 0)}
	} else {
		return &ParseContext{parent: parent, parser: parent.parser, tables: map[string]Table{}, warn: make([]string, 0)}
	}
}

func (ctx *ParseContext) addTables(tables ...Table) {
	for _, table := range tables {
		ctx.tables[table.Name()] = table
	}
}

func (ctx *ParseContext) getTable(name string) Table {
	t := ctx.tables[name]
	if t == nil && ctx.parent != nil {
		t = ctx.parent.getTable(name)
	}
	return t
}

func (ctx *ParseContext) getColumn(table string, column string) *Column {
	if table == "" {
		for _, t := range ctx.tables {
			if c := t.GetColumn(column); c != nil {
				return c
			}
		}
		if ctx.parent != nil {
			return ctx.parent.getColumn("", column)
		}
	} else {
		if t := ctx.getTable(table); t != nil {
			return t.GetColumn(column)
		}
	}
	return nil
}

func parseCreateTableStmt(stmt *ast.CreateTableStmt) *TableDefine {
	table := NewTableDefine(stmt.Table.Name.L)
	for _, def := range stmt.Cols {
		table.AddColumn(parseColumn(def))
	}

	return table
}

func (ctx *ParseContext) parseSelectStmt(as string, stmt *ast.SelectStmt) *TableDefine {
	nCtx := newParseContext(ctx)

	nCtx.parseWithClause(stmt.With)
	from := nCtx.parseJoin(UNNAMED, stmt.From.TableRefs)

	table := NewTableDefine(as)
	for _, field := range stmt.Fields.Fields {
		if field.WildCard != nil {
			if field.WildCard.Table.L == "" {
				table.Merge(from)
			} else {
				table.Merge(nCtx.getTable(field.WildCard.Table.L))
			}
		} else {
			col := ctx.parseExpr(field.Expr).as(getFieldExprName(field))
			if col == nil {
				ctx.warn = append(ctx.warn, "unhandled column "+field.OriginalText())
			} else {
				table.AddColumn(col)
			}
		}
	}

	ctx.warn = append(ctx.warn, nCtx.warn...)
	return table
}

func getFieldExprName(f *ast.SelectField) string {
	if f.AsName.L == "" {
		return f.OriginalText()
	} else {
		return f.AsName.L
	}
}

func booleanColumn(nullable bool) *Column {
	return NewColumn(UNNAMED, types.ETInt, nullable)
}

func mergeColumns(name string, columns ...*Column) *Column {
	// TODO
	return columns[0].as(name)
}

func (ctx *ParseContext) parseExpr(expr ast.ExprNode) *Column {
	if e, ok := expr.(*ast.ColumnNameExpr); ok {
		return ctx.getColumn(e.Name.Table.L, e.Name.Name.L)
	} else if _, ok := expr.(*ast.BetweenExpr); ok {
		return booleanColumn(false)
	} else if _, ok := expr.(*ast.IsTruthExpr); ok {
		return booleanColumn(false)
	} else if _, ok := expr.(*ast.IsNullExpr); ok {
		return booleanColumn(false)
	} else if b, ok := expr.(*ast.BinaryOperationExpr); ok {
		return mergeColumns(UNNAMED, ctx.parseExpr(b.L), ctx.parseExpr(b.R))
	} else if u, ok := expr.(*ast.UnaryOperationExpr); ok {
		return ctx.parseExpr(u.V)
	} else if c, ok := expr.(*ast.CaseExpr); ok {
		cases := utils.Map(c.WhenClauses, func(t *ast.WhenClause) *Column {
			return ctx.parseExpr(t.Expr)
		})
		return mergeColumns(UNNAMED, append(cases, ctx.parseExpr(c.Value))...)
	} else if _, ok := expr.(*ast.PatternRegexpExpr); ok {
		return booleanColumn(false)
	} else if _, ok := expr.(*ast.PatternInExpr); ok {
		return booleanColumn(false)
	} else if _, ok := expr.(*ast.PatternLikeExpr); ok {
		return booleanColumn(false)
	} else if _, ok := expr.(*ast.MatchAgainst); ok {
		return booleanColumn(false)
	} else if p, ok := expr.(*ast.ParenthesesExpr); ok {
		return ctx.parseExpr(p.Expr)
	} else if c, ok := expr.(*ast.FuncCastExpr); ok {
		return NewColumn(UNNAMED, c.Tp.EvalType(), false)
	} else if t, ok := expr.(ast.ValueExpr); ok {
		return NewColumn(UNNAMED, t.GetType().EvalType(), false)
	} else if f, ok := expr.(*ast.AggregateFuncExpr); ok {
		if tp := ctx.parser.FuncType(f.F); tp != nil {
			return &Column{*tp, UNNAMED}
		} else {
			ctx.warn = append(ctx.warn, "meet unknown aggregate func "+f.F)
			return NewColumn(UNNAMED, types.ETString, false)
		}
	} else if f, ok := expr.(*ast.FuncCallExpr); ok {
		if tp := ctx.parser.FuncType(f.FnName.L); tp != nil {
			return &Column{*tp, UNNAMED}
		} else {
			ctx.warn = append(ctx.warn, "meet unknown aggregate func "+f.FnName.O)
			return NewColumn(UNNAMED, types.ETString, false)
		}
	}
	panic(unimplementedNodeType(expr))
}

func (ctx *ParseContext) parseSubqueryExpr(as string, expr *ast.SubqueryExpr) Table {
	nCtx := newParseContext(ctx)
	t := nCtx.parseResultSetNode(as, expr.Query)
	ctx.warn = append(ctx.warn, nCtx.warn...)
	return t
}

func (ctx *ParseContext) parseWithClause(w *ast.WithClause) []Table {
	if w == nil {
		return []Table{}
	}
	return utils.Map(w.CTEs, func(t *ast.CommonTableExpression) Table {
		return ctx.parseSubqueryExpr(t.Name.L, t.Query)
	})
}

func (ctx *ParseContext) parseResultSetNode(as string, r ast.ResultSetNode) Table {
	if s, ok := r.(*ast.SelectStmt); ok {
		return ctx.parseSelectStmt(as, s)
	} else if s, ok := r.(*ast.SubqueryExpr); ok {
		return ctx.parseSubqueryExpr(as, s)
	} else if s, ok := r.(*ast.TableSource); ok {
		return NewTableRef(as, ctx.parseTableSource(s))
	} else if s, ok := r.(*ast.TableName); ok {
		return NewTableRef(as, ctx.getTable(s.Name.L))
	} else if s, ok := r.(*ast.Join); ok {
		return ctx.parseJoin(as, s)
	} else if s, ok := r.(*ast.SetOprStmt); ok {
		return ctx.parseSetOprStmt(as, s)
	} else {
		panic(unimplementedNodeType(r))
	}
}

func (ctx *ParseContext) parseJoin(as string, j *ast.Join) Table {
	left := ctx.parseResultSetNode(as, j.Left)
	if j.Right != nil {
		right := ctx.parseResultSetNode(UNNAMED, j.Right)
		return MergeTables(as, left, right)
	}
	return left
}

func (ctx *ParseContext) parseSetOprStmt(as string, s *ast.SetOprStmt) Table {
	cteTables := ctx.parseWithClause(s.With)
	ctx.addTables(cteTables...)

	cteTables = ctx.parseWithClause(s.SelectList.With)
	ctx.addTables(cteTables...)

	tables := make([]Table, len(s.SelectList.Selects))
	for i, node := range s.SelectList.Selects {
		if s, ok := node.(*ast.SelectStmt); ok {
			tables[i] = ctx.parseSelectStmt(UNNAMED, s)
		} else {
			panic("Unimplemented type " + reflect.TypeOf(node).Name() + ", in sql " + node.OriginalText())
		}
	}

	return MergeTables(as, tables...)
}

func (ctx *ParseContext) parseTableSource(ts *ast.TableSource) Table {
	if s, ok := ts.Source.(*ast.SelectStmt); ok {
		return ctx.parseSelectStmt(ts.AsName.L, s)
	}

	if s, ok := ts.Source.(*ast.SetOprStmt); ok {
		return ctx.parseSetOprStmt(ts.AsName.L, s)
	}

	if s, ok := ts.Source.(*ast.Join); ok {
		return ctx.parseJoin(ts.AsName.L, s)
	}

	if s, ok := ts.Source.(*ast.TableName); ok {
		return ctx.getTable(s.Name.L)
	}

	panic(unimplementedNodeType(ts.Source))
}

func unimplementedNodeType(node ast.Node) string {
	return "Unimplemented type " + reflect.TypeOf(node).String()
}
