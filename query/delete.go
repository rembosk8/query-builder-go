package query

import (
	"fmt"
)

type Delete struct {
	child
}

type DeleteCore struct {
	core

	table string
}

var _ Builder = &Delete{}

func (d *Delete) ToSQL() (sql string, err error) {
	// todo: check receiver type
	qb := qbInit(d)

	qb.buildDeleteFrom()
	qb.buildWhere()
	qb.buildReturning()

	return qb.strBuilder.String(), qb.err
}

func (d *Delete) ToSQLWithStmts() (sql string, args []any, err error) {
	qb := qbInit(d)

	qb.buildDeleteFrom()
	args = qb.buildWherePrepStmt(args)
	qb.buildReturning()

	return qb.SQLStmts(args)
}

func (d *Delete) Only() *Delete {
	o := Only{child{parent: d.parent}}
	d.parent = &o

	return d
}

func (d *Delete) Where(field string) *WherePart[*Delete] {
	// todo: check heap move
	w := Where{
		child: child{parent: d.parent},
		field: field,
	}

	d.parent = &w

	return &WherePart[*Delete]{
		b:     d,
		Where: &w,
	}
}

func (d *Delete) Returning(fields ...string) *Delete {
	if len(fields) == 0 {
		return d
	}
	r := Returning{
		child: child{parent: d.parent},
		rets:  fields,
	}

	d.parent = &r

	return d
}

func (qb *queryBuilder) buildDeleteFrom() {
	if qb.err != nil {
		return
	}

	upd := "DELETE FROM "
	if qb.only {
		upd += "ONLY "
	}
	_, qb.err = fmt.Fprint(qb.strBuilder, upd, qb.table)
}
