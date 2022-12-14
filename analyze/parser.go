package analyze

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	"github.com/pingcap/tidb/parser/types"
	"reflect"
	"runtime/debug"
	"strings"
	"tidb-sql-parser/utils"
)

type Parser struct {
	ctx    *ParseContext
	parser *parser.Parser

	FuncTypes       map[string]*Tp
	TransparentFunc map[string]bool
}

func NewParser() *Parser {
	p := &Parser{
		ctx:             newParseContext(nil),
		parser:          parser.New(),
		FuncTypes:       map[string]*Tp{},
		TransparentFunc: map[string]bool{},
	}
	p.ctx.parser = p
	return p
}

func (p *Parser) Warns() []string {
	return p.ctx.warn
}

func (p *Parser) DefineFunc(name string, tp *Tp) {
	p.FuncTypes[strings.ToLower(name)] = tp
}

func (p *Parser) DefineTransparentFunc(name string) {
	p.TransparentFunc[strings.ToLower(name)] = true
}

func (p *Parser) FuncType(name string) *Tp {
	return p.FuncTypes[strings.ToLower(name)]
}

func (p *Parser) TransparentFuncType(name string) bool {
	return p.TransparentFunc[strings.ToLower(name)] == true
}

func (p *Parser) AddDdl(sql string) {
	stmts, _, err := p.parser.ParseSQL(sql)
	if err != nil {
		panic(err)
	}
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.CreateTableStmt:
			p.ctx.addTables(parseCreateTableStmt(s))
		case *ast.CreateIndexStmt:
			continue
		default:
			panic(unimplementedNodeType(stmt))
		}
	}
}

func (p *Parser) Parse(sql string) (columns []*Column, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("stacktrace from error: \n" + string(debug.Stack()))
			columns = nil
			switch e := e.(type) {
			case error:
				err = e
			case string:
				err = &utils.RuntimeError{Message: e}
			default:
				err = &utils.RuntimeError{Message: "unknown error of type " + reflect.TypeOf(e).String()}
			}
		}
	}()

	p.ctx.warn = make([]string, 0)

	stmt, err := p.parser.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	if s, ok := stmt.(*ast.SelectStmt); ok {
		return p.ctx.parseSelectStmt(UNNAMED, s).Columns(), nil
	} else {
		panic(unimplementedNodeType(stmt))
	}
}

func (p *Parser) GetTable(table string) Table {
	return p.ctx.getTable(table)
}

type NormalizeDigestData struct {
	NormalizedSql string
	Digest        string
}

func (p *Parser) NormalizeDigest(sql string) *NormalizeDigestData {
	normalized, digest := parser.NormalizeDigest(sql)
	return &NormalizeDigestData{normalized, digest.String()}
}

type AstResult struct {
	Ast   []ast.StmtNode
	Warns []error
	Error error
}

func (p *Parser) GetAst(sql string) *AstResult {
	nodes, warns, err := p.parser.ParseSQL(sql)

	return &AstResult{nodes, warns, err}
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

func (ctx *ParseContext) addTable(table Table) {
	if table.Name() == UNNAMED {
		return
	}
	ctx.tables[table.Name()] = table
}

func (ctx *ParseContext) addTables(tables ...Table) {
	for _, table := range tables {
		ctx.addTable(table)
	}
}

func (ctx *ParseContext) getTable(name string) Table {
	t := ctx.tables[name]
	if t == nil && ctx.parent != nil {
		t = ctx.parent.getTable(name)
	}
	if t == nil {
		panic("failed to get table " + name)
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

	nCtx.parseWithClause(stmt.With, nCtx.addTable)

	var from Table
	if stmt.From != nil {
		from = nCtx.parseJoin(UNNAMED, stmt.From.TableRefs)
		nCtx.addTables(from)
	}

	table := NewTableDefine(as)
	for _, field := range stmt.Fields.Fields {
		if field.WildCard != nil {
			if field.WildCard.Table.L == "" {
				if from != nil {
					table.Merge(from)
				}
			} else {
				table.Merge(nCtx.getTable(field.WildCard.Table.L))
			}
		} else {
			raw := nCtx.parseExpr(field.Expr)
			col := raw.as(getFieldExprName(field, raw.Name))
			if col == nil {
				nCtx.warn = append(nCtx.warn, "unhandled column "+field.OriginalText())
			} else {
				table.AddColumn(col)
			}
		}
	}

	ctx.warn = append(ctx.warn, nCtx.warn...)
	return table
}

func getFieldExprName(f *ast.SelectField, fallback string) string {
	if f.AsName.O == "" {
		if fallback == UNNAMED {
			return f.OriginalText()
		} else {
			return fallback
		}
	} else {
		return f.AsName.O
	}
}

func booleanColumn(nullable bool) *Column {
	return NewColumn(UNNAMED, types.ETInt, nullable)
}

func mergeType(l, r types.EvalType) types.EvalType {
	if l < r {
		return r
	} else {
		return l
	}
}

func mergeTypes(tps ...*Tp) *Tp {
	if len(tps) == 0 {
		return nil
	}
	c := Tp{
		Type:     tps[0].Type,
		Nullable: tps[0].Nullable,
	}
	for i := 1; i < len(tps); i++ {
		c.Type = mergeType(c.Type, tps[i].Type)
		c.Nullable = c.Nullable || tps[i].Nullable
	}
	return tps[0]
}

func mergeColumns(name string, columns ...*Column) *Column {
	tps := utils.Map(columns, func(c *Column) *Tp {
		return &c.Tp
	})
	tp := mergeTypes(tps...)
	return NewColumn(name, tp.Type, tp.Nullable)
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
		if b.Op == opcode.Div {
			return NewColumn(UNNAMED, types.ETReal, false)
		}
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
		tp := ctx.parseFuncType(f.F, f.Args)
		return &Column{*tp, UNNAMED, ""}
	} else if f, ok := expr.(*ast.FuncCallExpr); ok {
		tp := ctx.parseFuncType(f.FnName.O, f.Args)
		return &Column{*tp, UNNAMED, ""}
	} else if f, ok := expr.(*ast.WindowFuncExpr); ok {
		tp := ctx.parseFuncType(f.F, f.Args)
		return &Column{*tp, UNNAMED, ""}
	}
	panic(unimplementedNodeType(expr))
}

func (ctx *ParseContext) parseFuncType(name string, args []ast.ExprNode) *Tp {
	if tp := ctx.parser.FuncType(name); tp != nil {
		return tp
	}
	if ctx.parser.TransparentFuncType(name) {
		args := utils.Map(args, func(t ast.ExprNode) *Tp {
			return &ctx.parseExpr(t).Tp
		})
		return args[0]
	}
	ctx.warn = append(ctx.warn, "Unknown type of func "+name+". Use Parser.DefineFunc or Parser.DefineTransparentFunc to provide accurate type.")
	return &Tp{
		types.ETString,
		false,
	}
}

func (ctx *ParseContext) parseSubqueryExpr(as string, expr *ast.SubqueryExpr) Table {
	nCtx := newParseContext(ctx)
	t := nCtx.parseResultSetNode(as, expr.Query)
	ctx.warn = append(ctx.warn, nCtx.warn...)
	return t
}

func (ctx *ParseContext) parseWithClause(w *ast.WithClause, cb func(Table)) {
	if w == nil {
		return
	}
	for _, c := range w.CTEs {
		cb(ctx.parseSubqueryExpr(c.Name.L, c.Query))
	}
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
	ctx.addTable(left)
	if j.Right != nil {
		right := ctx.parseResultSetNode(UNNAMED, j.Right)
		ctx.addTable(right)
		return MergeTables(as, left, right)
	}
	return left
}

func (ctx *ParseContext) parseSetOprStmt(as string, s *ast.SetOprStmt) Table {
	ctx.parseWithClause(s.With, ctx.addTable)

	ctx.parseWithClause(s.SelectList.With, ctx.addTable)

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
		return NewTableRef(ts.AsName.L, ctx.getTable(s.Name.L))
	}

	panic(unimplementedNodeType(ts.Source))
}

func unimplementedNodeType(node ast.Node) string {
	return "Unimplemented type " + reflect.TypeOf(node).String()
}
