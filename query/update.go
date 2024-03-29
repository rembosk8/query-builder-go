package query

import (
	"fmt"
	"strings"
)

type UpdateCore struct {
	core

	table string
}

type setValue struct {
	child

	// todo: impl FiledValue struct and make it possible to set many at one call
	//  or, may be one entity one value could be better, need to check.
	fvs []FiledValue
}

func (sv *setValue) Set(fvs []any) {
	for i := 0; i < len(fvs); i++ {
		switch t := fvs[i].(type) {
		case string:
			sv.fvs = append(sv.fvs, FldVal(t, fvs[i+1]))
			i++
		case FiledValue:
			sv.fvs = append(sv.fvs, t)
		default:
			panic(fmt.Sprintf("unknown key type in Sets %T", t))
		}
	}
}

type Update struct {
	child
}

func (u Update) Set(field string, value any) *Update {
	usv := setValue{
		child: child{parent: u.parent},
		fvs: []FiledValue{{
			field: field,
			value: value,
		}},
	}

	return &Update{child{parent: &usv}}
}

func (u Update) Sets(filedValues ...any) *Update {
	if len(filedValues) == 0 {
		return &u
	}

	usv := setValue{
		child: child{parent: u.parent},
		fvs:   []FiledValue{},
	}

	usv.Set(filedValues)

	return &Update{child{parent: &usv}}
}

// todo: move to a separate file.

func FldVal[T any](field string, value T) FiledValue {
	return FiledValue{
		field: field,
		value: value,
	}
}

type FiledValue struct {
	field string
	value any
}

var _ Builder = &Update{}

func (u Update) ToSQL() (sql string, err error) {
	qb := qbInit(u)
	u.Parent()

	qb.buildUpdateTable()
	qb.buildSet()
	qb.buildWhere()
	qb.buildReturning()

	return qb.strBuilder.String(), qb.err
}

var _ parenter = Update{}

func (u Update) ToSQLWithStmts() (sql string, args []any, err error) {
	qb := qbInit(u)
	args = qb.buildPrepStatement()

	return qb.strBuilder.String(), args, nil
}

func (u Update) Where(field string) *WherePart[*Update] {
	w := Where{
		child: child{parent: u.parent},
		field: field,
	}

	u.parent = &w

	return &WherePart[*Update]{
		b:     &u,
		Where: &w,
	}
}

type Only struct {
	child
}

func (u Update) Only() *Update {
	only := Only{child{parent: u.parent}}
	u.parent = &only
	return &u
}

type Returning struct {
	child

	rets []string
}

func (u Update) Returning(fields ...string) *Update {
	if len(fields) == 0 {
		return &u
	}
	r := Returning{
		child: child{parent: u.parent},
		rets:  fields,
	}

	u.parent = &r

	return &u
}

func (qb *queryBuilder) buildPrepStatement() (args []any) {
	qb.buildUpdateTable()
	args = qb.buildSetStmt()
	args = qb.buildWherePrepStmt(args)
	qb.buildReturning()
	return
}

func (qb *queryBuilder) buildUpdateTable() {
	if qb.err != nil {
		return
	}

	upd := "UPDATE "
	if qb.only {
		upd += "ONLY "
	}
	_, qb.err = fmt.Fprint(qb.strBuilder, upd, qb.table)
}

func (qb *queryBuilder) buildSet() {
	if qb.err != nil {
		return
	}

	if len(qb.fields) == 0 {
		qb.err = ErrUpdateValuesNotSet
		return
	}

	//todo: try without format.
	_, qb.err = fmt.Fprintf(qb.strBuilder, " SET %s = %s", qb.fields[0], qb.indentBuilder.Value(qb.values[0]))

	for i := 1; i < len(qb.fields); i++ {
		_, qb.err = fmt.Fprintf(qb.strBuilder, ", %s = %s", qb.fields[i], qb.indentBuilder.Value(qb.values[i]))
	}
}

func (qb *queryBuilder) buildSetStmt() (args []any) {
	if qb.err != nil {
		return
	}
	if len(qb.fields) == 0 {
		qb.err = ErrUpdateValuesNotSet
		return nil
	}
	args = make([]any, 0, len(qb.fields))

	_, qb.err = fmt.Fprint(qb.strBuilder, " SET ")

	setStmt := func(i int) {
		if qb.indentBuilder.IsStandard(qb.values[i]) {
			_, qb.err = fmt.Fprintf(qb.strBuilder, "%s = %s", qb.fields[i], qb.values[i])
			return
		}
		_, qb.err = fmt.Fprintf(qb.strBuilder, "%s = $%d", qb.fields[i], len(args)+1)
		args = append(args, qb.values[i])
	}

	setStmt(0)
	for i := 1; i < len(qb.fields); i++ {
		_, qb.err = fmt.Fprint(qb.strBuilder, ", ")
		setStmt(i)
	}

	return args
}

func (qb *queryBuilder) buildReturning() {
	if len(qb.returning) == 0 {
		return
	}
	if qb.err != nil {
		return
	}

	_, qb.err = fmt.Fprint(
		qb.strBuilder,
		" RETURNING ",
	)

	_, qb.err = fmt.Fprint(
		qb.strBuilder,
		strings.Join(qb.returning, ", "),
	)
}
