package query

import (
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/internal/helpers/pointer"
	identity2 "github.com/rembosk8/query-builder-go/internal/identity"
)

type Builder interface {
	ToSQL() (sql string, err error)
	ToSQLWithStmts() (sql string, args []any, err error)
}

type baseQuery struct {
	// todo :check what is needed.

	indentBuilder *identity2.Builder
	tag           string
}

type queryBuilder struct {
	cols     []identity2.Identity
	table    identity2.Identity
	wheres   []*Where //todo: [][]Where or []OrWhere -> OrWhere{[]AndWhere}. idea: to have a map of []Where, where key is OR-group name
	offset   *uint
	limit    *uint
	orderBys []*Order
	joins    []*Join

	// update.
	only      bool
	returning []identity2.Identity //todo: add to insert

	// insert or update.
	fields []identity2.Identity
	values []any

	err           error
	indentBuilder *identity2.Builder
	strBuilder    *strings.Builder
	//todo: try to add arena https://uptrace.dev/blog/golang-memory-arena.html
}

func (qb *queryBuilder) SQLStmts(args []any) (sql string, argsOut []any, err error) {
	if qb.err != nil {
		return "", nil, qb.err
	}

	return qb.strBuilder.String(), args, nil
}

func (qb *queryBuilder) SQL() (sql string, err error) {
	if qb.err != nil {
		return "", qb.err
	}

	return qb.strBuilder.String(), nil
}

func (qb *queryBuilder) collect(p any) { //nolint:cyclop
	par, ok := p.(parenter)
	if ok {
		qb.collect(par.Parent())
	}

	if par == nil {
		return
	}

	switch q := p.(type) {
	case Select, Update, *Delete, Insert:
		return
	case *SelectCore:
		qb.indentBuilder = q.indentBuilder

		if len(q.fields) > 0 {
			qb.cols = qb.indentBuilder.Idents(q.fields...)
		}
		qb.table = qb.indentBuilder.Ident(q.table)
	case *Join:
		qb.joins = append(qb.joins, q)
	case *UpdateCore:
		qb.indentBuilder = q.indentBuilder
		qb.table = qb.indentBuilder.Ident(q.table)
	case *DeleteCore:
		qb.indentBuilder = q.indentBuilder
		qb.table = qb.indentBuilder.Ident(q.table)
	case *InsertCore:
		qb.indentBuilder = q.indentBuilder
		qb.table = qb.indentBuilder.Ident(q.table)
	case *Where:
		qb.wheres = append(qb.wheres, q)
	case *Offset:
		qb.offset = pointer.To(q.offset)
	case *Limit:
		qb.limit = pointer.To(q.limit)
	case *Order:
		qb.orderBys = append(qb.orderBys, q)
	case *Only:
		qb.only = true
	case *setValue:
		for i := range q.fvs {
			qb.fields = append(qb.fields, qb.indentBuilder.Ident(q.fvs[i].field))
			qb.values = append(qb.values, q.fvs[i].value)
		}
	case *Returning:
		qb.returning = append(qb.returning, qb.indentBuilder.Idents(q.rets...)...)
	default:
		panic(fmt.Sprintf("wrong type in collect %T", p))
	}
}

func (qb *queryBuilder) getFields() string {
	if len(qb.cols) == 0 {
		return "*"
	}

	return strings.Join(qb.cols, ", ")
}

func (qb *queryBuilder) buildSelectFrom() {
	if qb.err != nil {
		return
	}

	_, qb.err = fmt.Fprintf(qb.strBuilder, "SELECT %s FROM %s", qb.getFields(), qb.table)
	for _, j := range qb.joins {
		if qb.err != nil {
			return
		}
		_, qb.err = fmt.Fprint(qb.strBuilder, j.String(qb.indentBuilder))
	}
}

func (qb *queryBuilder) buildWhere() {
	if len(qb.wheres) == 0 {
		return
	}

	_, qb.err = fmt.Fprint(qb.strBuilder, " WHERE ")
	if qb.err != nil {
		return
	}

	_, qb.err = fmt.Fprint(qb.strBuilder, qb.wheres[0].String(qb.indentBuilder))
	if qb.err != nil {
		return
	}

	for i := 1; i < len(qb.wheres); i++ {
		_, qb.err = fmt.Fprint(qb.strBuilder, " AND ", qb.wheres[i].String(qb.indentBuilder))
		if qb.err != nil {
			return
		}
	}
}

func (qb *queryBuilder) buildWherePrepStmt(args []any) []any {
	if qb.err != nil {
		return nil
	}
	if len(qb.wheres) == 0 {
		return args
	}
	var vals []any
	cnt := len(args) + 1

	_, qb.err = fmt.Fprint(qb.strBuilder, " WHERE ")
	vals, qb.err = qb.wheres[0].PrepStmtString(cnt, qb.strBuilder, qb.indentBuilder)
	if qb.err != nil {
		return nil
	}
	args = append(args, vals...)

	cnt += len(vals)
	for i := 1; i < len(qb.wheres); i++ {
		_, qb.err = fmt.Fprint(qb.strBuilder, " AND ")
		if qb.err != nil {
			return nil
		}
		vals, qb.err = qb.wheres[i].PrepStmtString(cnt, qb.strBuilder, qb.indentBuilder)
		if qb.err != nil {
			return nil
		}
		args = append(args, vals...)
		cnt += len(vals)
	}

	return args
}

func (qb *queryBuilder) BuildOffset() {
	if qb.err != nil {
		return
	}
	if qb.offset == nil {
		return
	}

	_, qb.err = fmt.Fprintf(qb.strBuilder, " OFFSET %d", *qb.offset)
}

func (qb *queryBuilder) BuildLimit() {
	if qb.err != nil {
		return
	}
	if qb.limit == nil {
		return
	}
	_, qb.err = fmt.Fprintf(qb.strBuilder, " LIMIT %d", *qb.limit)
}

func (qb *queryBuilder) BuildOrderBy() {
	if len(qb.orderBys) == 0 || qb.err != nil {
		return
	}
	_, qb.err = fmt.Fprint(qb.strBuilder, " ORDER BY ", qb.orderBys[0].String(qb.indentBuilder))
	if qb.err != nil {
		return
	}

	for i := 1; i < len(qb.orderBys); i++ {
		_, qb.err = fmt.Fprint(qb.strBuilder, ", ", qb.orderBys[i].String(qb.indentBuilder))
		if qb.err != nil {
			return
		}
	}
}
