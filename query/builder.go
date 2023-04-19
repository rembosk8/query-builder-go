package query

import (
	"errors"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/rembosk8/query-builder-go/query/indent"
	"github.com/rembosk8/query-builder-go/stringer"
)

const all = "*"

var ErrTableNotSet = errors.New("table name not provided")

type Builder struct {
	fields        []indent.Indent // select <fields>
	table         *indent.Indent  // from <table>
	wheres        []Where
	indentBuilder *indent.Builder
	offset        *uint
	limit         *uint
	orderBys      []Order
}

func New(indBuilder *indent.Builder) Builder {
	if indBuilder == nil {
		indBuilder = indent.NewBuilder()
	}

	return Builder{indentBuilder: indBuilder}
}

func (b Builder) Build() (sql string, args []any, err error) {
	if b.table == nil {
		return "", nil, ErrTableNotSet
	}
	sql, args = b.buildPrepStatement()

	return sql, args, nil
}

func (b Builder) BuildPlain() (sql string, err error) {
	if b.table == nil {
		return "", ErrTableNotSet
	}
	return b.buildSqlPlain(), nil
}

func (b Builder) From(tableName string) Builder {
	b.table = pointer.To(b.indentBuilder.Indent(tableName))
	return b
}

func (b Builder) Select(fields ...string) Builder {
	for _, f := range fields {
		b.fields = append(b.fields, b.indentBuilder.Indent(f))
	}
	return b
}

func (b Builder) getFields() string {
	if len(b.fields) == 0 {
		return all
	}

	return stringer.Join(b.fields, ", ")
}

func (b Builder) buildSelectFrom() string {
	return fmt.Sprintf("SELECT %s FROM %s", b.getFields(), b.table.String())
}

func (b Builder) buildWherePlain() string {
	if len(b.wheres) == 0 {
		return ""
	}
	return fmt.Sprintf(" WHERE %s", stringer.Join(b.wheres, " AND ")) //todo: build AND and OR separately
}

func (b Builder) buildOffset() string {
	if b.offset == nil {
		return ""
	}

	return fmt.Sprintf(" OFFSET %d", *b.offset)
}

func (b Builder) buildLimit() string {
	if b.limit == nil {
		return ""
	}

	return fmt.Sprintf(" LIMIT %d", *b.limit)
}

func (b Builder) buildWherePrepStmt() (string, []any) {
	if len(b.wheres) == 0 {
		return "", nil
	}
	var args []any
	cnt := 1
	prepStmt, vals := b.wheres[0].PrepStmtString(cnt)
	args = append(args, vals...)
	sql := fmt.Sprintf(" WHERE %s", prepStmt)

	cnt += len(vals)
	for i := 1; i < len(b.wheres); i++ {
		prepStmt, vals = b.wheres[i].PrepStmtString(cnt)
		args = append(args, vals...)
		sql += fmt.Sprintf(" AND %s", prepStmt)
		cnt += len(vals)
	}

	return sql, args
}

func (b Builder) buildSqlPlain() string {
	sql := b.buildSelectFrom()
	if len(b.wheres) > 0 {
		sql += b.buildWherePlain()
	}
	sql += b.buildOrderBy()
	sql += b.buildOffset()
	sql += b.buildLimit()

	return sql
}

func (b Builder) buildPrepStatement() (sql string, args []any) {
	sql = b.buildSelectFrom()
	if len(b.wheres) > 0 {
		var sqlWhere string
		sqlWhere, args = b.buildWherePrepStmt()
		sql += sqlWhere
	}
	sql += b.buildOrderBy()
	sql += b.buildOffset()
	sql += b.buildLimit()

	return sql, args
}

func (b Builder) Where(columnName string) wherePart {
	return wherePart{
		column: b.indentBuilder.Indent(columnName),
		b:      b,
	}
}

func (b Builder) Offset(n uint) Builder {
	b.offset = pointer.To(n)
	return b
}

func (b Builder) Limit(n uint) Builder {
	b.limit = pointer.To(n)
	return b
}

func (b Builder) OrderBy(fieldName string) orderPart {
	return orderPart{
		column: b.indentBuilder.Indent(fieldName),
		b:      b,
	}
}

func (b Builder) buildOrderBy() string {
	if len(b.orderBys) == 0 {
		return ""
	}

	return fmt.Sprintf(" ORDER BY %s", stringer.Join(b.orderBys, ", "))
}
