package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/identity"
)

var _ fmt.Stringer = &filedValue{}

type filedValue struct {
	field identity.Identity
	value identity.Value
}

func (f filedValue) String() string {
	return fmt.Sprintf("%s = %s", f.field.String(), f.value.String())
}

func (f filedValue) StringStmt(i uint16) (sql string, v any) {
	if f.value.IsStandard() {
		return f.String(), nil
	}
	return fmt.Sprintf("%s = $%d", f.field.String(), i), f.value.Value
}

type Update struct {
	baseQuery

	fieldValue []filedValue
	returning  []identity.Identity
	only       bool
}

var _ sqler = &Update{}

func (u Update) ToSql() (sql string, err error) {
	if err = u.initBuild(); err != nil {
		return "", err
	}
	u.buildSqlPlain()

	return u.strBuilder.String(), nil
}

func (u Update) ToSqlWithStmts() (sql string, args []any, err error) {
	if err = u.initBuild(); err != nil {
		return "", nil, err
	}
	args = u.buildPrepStatement()

	return u.strBuilder.String(), args, nil
}

func (u Update) Set(field string, value any) Update {
	u.fieldValue = append(u.fieldValue, filedValue{
		field: u.indend(field),
		value: u.value(value),
	})
	return u
}

func (u Update) Where(columnName string) wherePart[*Update] {
	return wherePart[*Update]{
		column: u.indend(columnName),
		b:      &u,
	}
}

func (u Update) Only() Update {
	u.only = true
	return u
}

func (u Update) Returning(fields ...string) Update {
	for _, f := range fields {
		u.returning = append(u.returning, u.indend(f))
	}

	return u
}

func (u *Update) buildSqlPlain() {
	u.buildUpdateTable()
	u.buildSet()
	u.buildWhere()
	u.buildReturning()
}

func (u *Update) buildPrepStatement() (args []any) {
	u.buildUpdateTable()
	args = u.buildSetStmt()
	args = u.buildWherePrepStmt(args)
	u.buildReturning()
	return
}

func (u *Update) buildUpdateTable() {
	if u.err != nil {
		return
	}

	upd := "UPDATE "
	if u.only {
		upd += "ONLY "
	}
	_, u.err = fmt.Fprint(u.strBuilder, upd, u.table.String())
}

func (u *Update) buildSet() {
	if u.err != nil {
		return
	}

	if len(u.fieldValue) == 0 {
		u.err = ErrUpdateValuesNotSet
		return
	}

	_, u.err = fmt.Fprintf(u.strBuilder, " SET %s", stringer.Join(u.fieldValue, ", "))
}

func (u *Update) buildSetStmt() (args []any) {
	if u.err != nil {
		return
	}

	if len(u.fieldValue) == 0 {
		u.err = ErrUpdateValuesNotSet
		return nil
	}

	var num uint16 = 1
	args = make([]any, 0, len(u.fieldValue))

	sql, v := u.fieldValue[0].StringStmt(num)
	if v != nil {
		args = append(args, v)
		num++
	}

	_, u.err = fmt.Fprintf(u.strBuilder, " SET %s", sql)

	for i := 1; i < len(u.fieldValue); i++ {

		sql, v = u.fieldValue[i].StringStmt(num)
		if v != nil {
			args = append(args, v)
			num++
		}
		_, u.err = fmt.Fprintf(u.strBuilder, ", %s", sql)
	}

	return args
}

func (u *Update) buildReturning() {
	if u.err != nil {
		return
	}

	if len(u.returning) == 0 {
		return
	}

	_, u.err = fmt.Fprintf(u.strBuilder, " RETURNING %s", stringer.Join(u.returning, ", "))
}
