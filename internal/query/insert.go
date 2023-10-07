package query

import (
	"fmt"
	"strings"
)

type InsertCore struct {
	core

	table string
}

type Insert struct {
	child
}

var _ Builder = &Insert{}

func (i Insert) ToSQL() (sql string, err error) {
	qb := qbInit(i)
	qb.buildInsertInto()
	qb.buildValues()

	return qb.SQL()
}

func (i Insert) ToSQLWithStmts() (sql string, args []any, err error) {
	qb := qbInit(i)
	qb.buildInsertInto()
	args = qb.buildValueStmts()

	return qb.SQLStmts(args)
}

func (i Insert) Set(field string, value any) *Insert {
	sv := setValue{
		child: child{parent: i.parent},
		fvs: []FiledValue{{
			field: field,
			value: value,
		}},
	}

	i.parent = &sv

	return &i
}

func (qb *queryBuilder) buildInsertInto() {
	if qb.err != nil {
		return
	}

	_, qb.err = fmt.Fprint(qb.strBuilder, "INSERT INTO ", qb.table)
}

func (qb *queryBuilder) buildDefaultValues() {
	_, qb.err = fmt.Fprintf(
		qb.strBuilder,
		" DEFAULT VALUES",
	)
}

func (qb *queryBuilder) buildValues() {
	if qb.err != nil {
		return
	}

	if len(qb.fields) == 0 {
		qb.buildDefaultValues()
		return
	}

	_, qb.err = fmt.Fprintf(
		qb.strBuilder,
		" (%s) VALUES (%s)",
		strings.Join(qb.fields, ", "),
		strings.Join(qb.indentBuilder.Values(qb.values), ", "), //todo: check if it's possible to not use Values here
	)
}

func (qb *queryBuilder) buildValueStmts() (args []any) {
	if qb.err != nil {
		return
	}

	if len(qb.fields) == 0 {
		qb.buildDefaultValues()
		return
	}

	_, qb.err = fmt.Fprintf(
		qb.strBuilder,
		" (%s) VALUES (%s)",
		strings.Join(qb.fields, ", "),
		genNums(1, len(qb.values)),
	)

	return qb.values
}
