package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/helpers/pointer"
	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/indent"
)

const all = "*"

var ErrTableNotSet = errors.New("table name not provided")

type Builder struct {
	fields   []indent.Indent // select <fields>
	table    *indent.Indent  // from <table>
	wheres   []*Where
	offset   *uint
	limit    *uint
	orderBys []Order

	err error

	indentBuilder *indent.Builder
	strBuilder    strings.Builder
	//l int  //precalculated len of statement
}

func New(indBuilder *indent.Builder) Builder {
	if indBuilder == nil {
		indBuilder = indent.NewBuilder()
	}

	return Builder{indentBuilder: indBuilder}
}

func (b Builder) Build() (sql string, args []any, err error) {
	if b.err != nil {
		return "", nil, err
	}
	if b.table == nil {
		return "", nil, ErrTableNotSet
	}
	args = b.buildPrepStatement()

	return b.strBuilder.String(), args, nil
}

func (b Builder) BuildPlain() (sql string, err error) {
	if b.err != nil {
		return "", err
	}
	if b.table == nil {
		return "", ErrTableNotSet
	}

	b.buildSqlPlain()
	if b.err != nil {
		return "", err
	}

	return b.strBuilder.String(), nil
}

func (b Builder) From(tableName string) Builder {
	b.table = pointer.To(b.indentBuilder.Indent(tableName))
	return b
}

func (b Builder) Select(fields ...string) Builder {
	//todo: create a separate method which can extract fields tags: by fields, by the whole struct
	for _, f := range fields {
		b.fields = append(b.fields, b.indentBuilder.Indent(f))
	}
	return b
}

func (b *Builder) getFields() string {
	if len(b.fields) == 0 {
		return all
	}

	return stringer.Join(b.fields, ", ")
}

func (b *Builder) buildSelectFrom() {
	if b.err != nil {
		return
	}

	_, b.err = fmt.Fprintf(&b.strBuilder, "SELECT %s FROM %s", b.getFields(), b.table.String())
}

func (b *Builder) buildWherePlain() {
	if b.err != nil {
		return
	}
	if len(b.wheres) == 0 {
		return
	}
	_, b.err = fmt.Fprintf(&b.strBuilder, " WHERE %s", stringer.Join(b.wheres, " AND ")) //todo: build AND and OR separately
}

func (b *Builder) buildOffset() {
	if b.err != nil {
		return
	}
	if b.offset == nil {
		return
	}

	_, b.err = fmt.Fprintf(&b.strBuilder, " OFFSET %d", *b.offset)
}

func (b *Builder) buildLimit() {
	if b.err != nil {
		return
	}
	if b.limit == nil {
		return
	}
	_, b.err = fmt.Fprintf(&b.strBuilder, " LIMIT %d", *b.limit)
}

func (b *Builder) buildWherePrepStmt() []any {
	if len(b.wheres) == 0 {
		return nil
	}
	if b.err != nil {
		return nil
	}
	var args, vals []any
	cnt := 1

	_, b.err = fmt.Fprint(&b.strBuilder, " WHERE ")
	vals, b.err = b.wheres[0].PrepStmtString(cnt, &b.strBuilder)
	if b.err != nil {
		return nil
	}
	args = append(args, vals...)

	cnt += len(vals)
	for i := 1; i < len(b.wheres); i++ {
		_, b.err = fmt.Fprint(&b.strBuilder, " AND ")
		if b.err != nil {
			return nil
		}
		vals, b.err = b.wheres[i].PrepStmtString(cnt, &b.strBuilder)
		if b.err != nil {
			return nil
		}
		args = append(args, vals...)
		cnt += len(vals)
	}

	return args
}

func (b *Builder) buildSqlPlain() {
	b.buildSelectFrom()

	b.buildWherePlain()

	b.buildOrderBy()
	b.buildOffset()
	b.buildLimit()
}

func (b *Builder) buildPrepStatement() (args []any) {
	b.buildSelectFrom()

	args = b.buildWherePrepStmt()

	b.buildOrderBy()
	b.buildOffset()
	b.buildLimit()

	return args
}

func (b Builder) Where(columnName string) wherePart {
	//b.l += len(columnName)
	return wherePart{
		column: b.indentBuilder.Indent(columnName),
		b:      b,
	}
}

func (b Builder) Offset(n uint) Builder {
	//b.l += 12
	b.offset = pointer.To(n)
	return b
}

func (b Builder) Limit(n uint) Builder {
	//b.l += 12
	b.limit = pointer.To(n)
	return b
}

func (b Builder) OrderBy(fieldName string) orderPart {
	//b.l += 15 + len(fieldName)
	return orderPart{
		column: b.indentBuilder.Indent(fieldName),
		b:      b,
	}
}

func (b *Builder) buildOrderBy() {
	if b.err != nil {
		return
	}
	if len(b.orderBys) == 0 {
		return
	}
	_, b.err = fmt.Fprintf(&b.strBuilder, " ORDER BY %s", stringer.Join(b.orderBys, ", "))
}
