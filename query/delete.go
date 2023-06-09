package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/identity"
)

type Delete struct {
	baseQuery
	only      bool
	returning []identity.Identity
}

var _ sqler = &Delete{}

func (d Delete) ToSql() (sql string, err error) {
	if err = d.initBuild(); err != nil {
		return "", err
	}
	d.buildSqlPlain()

	return d.strBuilder.String(), nil
}

func (d Delete) ToSqlWithStmts() (sql string, args []any, err error) {
	if err = d.initBuild(); err != nil {
		return "", nil, err
	}
	args = d.buildPrepStatement()

	return d.strBuilder.String(), args, nil
}

func (d Delete) Only() Delete {
	d.only = true
	return d
}

func (d Delete) Where(field string) wherePart[*Delete] {
	return wherePart[*Delete]{
		column: d.indend(field),
		b:      &d,
	}
}

func (d Delete) Returning(fields ...string) Delete {
	for _, f := range fields {
		d.returning = append(d.returning, d.indend(f))
	}

	return d
}

func (d *Delete) buildDeleteFrom() {
	if d.err != nil {
		return
	}

	upd := "DELETE FROM "
	if d.only {
		upd += "ONLY "
	}
	_, d.err = fmt.Fprint(d.strBuilder, upd, d.table.String())
}

func (d Delete) buildSqlPlain() {
	d.buildDeleteFrom()
	d.buildWhere()
	d.buildReturning()
}

func (d Delete) buildPrepStatement() (args []any) {
	d.buildDeleteFrom()
	args = d.buildWherePrepStmt(args)
	d.buildReturning()

	return
}

func (d *Delete) buildReturning() {
	if d.err != nil {
		return
	}

	if len(d.returning) == 0 {
		return
	}

	_, d.err = fmt.Fprintf(d.strBuilder, " RETURNING %s", stringer.Join(d.returning, ", "))
}
