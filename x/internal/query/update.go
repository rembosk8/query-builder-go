package query

//
//import (
//	"fmt"
//	"strings"
//
//	"github.com/rembosk8/query-builder-go/x/internal/identity"
//)
//
//var _ fmt.Stringer = &filedValue{}
//
//type filedValue struct {
//	field identity.Identity
//	value identity.Value
//}
//
//func (f filedValue) String(idb *identity.Builder) string {
//	return fmt.Sprintf("%s = %s", idb.Ident(f.field), idb.Value(f.value))
//}
//
//func (f filedValue) StringStmt(idb *identity.Builder, i uint16) (sql string, v any) {
//	if idb.IsStandard(f.value) {
//		return f.String(idb), nil
//	}
//	return fmt.Sprintf("%s = $%d", idb.Ident(f.field), i), f.value
//}
//
//type Update struct {
//	baseQuery
//
//	fieldValue []filedValue
//	returning  []string
//	only       bool
//}
//
//var _ Builder = &Update{}
//
//func (u Update) ToSQL() (sql string, err error) {
//	if err := u.initBuild(); err != nil {
//		return "", err
//	}
//	u.buildSQLPlain()
//
//	return u.strBuilder.String(), nil
//}
//
//func (u Update) ToSQLWithStmts() (sql string, args []any, err error) {
//	if err := u.initBuild(); err != nil {
//		return "", nil, err
//	}
//	args = u.buildPrepStatement()
//
//	return u.strBuilder.String(), args, nil
//}
//
//func (u Update) Set(field string, value any) Update {
//	u.fieldValue = append(u.fieldValue, filedValue{
//		field: u.ident(field),
//		value: u.value(value),
//	})
//	return u
//}
//
//func (u Update) Where(columnName string) WherePart[*Update] { //nolint:revive
//	return WherePart[*Update]{
//		Where: Where{
//			child: child{parent: nil}, //todo: set u
//			field: columnName,
//		},
//		b: &u,
//	}
//}
//
//func (u Update) Only() Update {
//	u.only = true
//	return u
//}
//
//func (u Update) Returning(fields ...string) Update {
//	for _, f := range fields {
//		u.returning = append(u.returning, u.ident(f))
//	}
//
//	return u
//}
//
//func (u *Update) buildSQLPlain() {
//	//u.buildUpdateTable()
//	//u.buildSet()
//	//u.buildWhere()
//	//u.buildReturning()
//}
//
//func (u *Update) buildPrepStatement() (args []any) {
//	//u.buildUpdateTable()
//	//args = u.buildSetStmt()
//	//args = u.buildWherePrepStmt(args)
//	//u.buildReturning()
//	return
//}
//
//func (u *Update) buildUpdateTable(idb *identity.Builder) {
//	if u.err != nil {
//		return
//	}
//
//	//upd := "UPDATE "
//	//if u.only {
//	//	upd += "ONLY "
//	//}
//	//_, u.err = fmt.Fprint(u.strBuilder, upd, idb.Ident(u.table))
//}
//
//func (u *Update) buildSet() {
//	//if u.err != nil {
//	//	return
//	//}
//	//
//	//if len(u.fieldValue) == 0 {
//	//	u.err = ErrUpdateValuesNotSet
//	//	return
//	//}
//	//
//	//_, u.err = fmt.Fprintf(u.strBuilder, " SET %s", stringer.Join(u.fieldValue, ", "))
//}
//
//func (u *Update) buildSetStmt() (args []any) {
//	if u.err != nil {
//		return
//	}
//
//	//if len(u.fieldValue) == 0 {
//	//	u.err = ErrUpdateValuesNotSet
//	//	return nil
//	//}
//	//
//	//var num uint16 = 1
//	//args = make([]any, 0, len(u.fieldValue))
//	//
//	//sql, v := u.fieldValue[0].StringStmt(num)
//	//if v != nil {
//	//	args = append(args, v)
//	//	num++
//	//}
//	//
//	//_, u.err = fmt.Fprintf(u.strBuilder, " SET %s", sql)
//	//
//	//for i := 1; i < len(u.fieldValue); i++ {
//	//	sql, v = u.fieldValue[i].StringStmt(num)
//	//	if v != nil {
//	//		args = append(args, v)
//	//		num++
//	//	}
//	//	_, u.err = fmt.Fprintf(u.strBuilder, ", %s", sql)
//	//}
//
//	return args
//}
//
//func (u *Update) buildReturning(idb *identity.Builder) {
//	if u.err != nil {
//		return
//	}
//
//	if len(u.returning) == 0 {
//		return
//	}
//
//	_, u.err = fmt.Fprintf(u.strBuilder, " RETURNING %s", strings.Join(idb.Idents(u.returning...), ", "))
//}
