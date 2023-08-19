package query

import (
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/x/internal/identity"
)

type selecter interface {
	Builder

	Where(string) *WherePart[*Select]
}

type parenter interface {
	Parent() any
}

type child struct {
	parent parenter
}

func (c *child) Parent() any {
	if c == nil {
		return nil
	}

	return c.parent
}

type Select struct {
	child
}

func (s Select) Where(field string) *WherePart[*Select] {
	w := Where{
		child: child{parent: s.parent},
		field: field,
	}

	s.parent = &w

	return &WherePart[*Select]{
		Where: &w,
		b:     &s,
	}
}

var _ selecter = &Select{}

type core struct {
	child
	indentBuilder *identity.Builder
}

type SelectCore struct {
	core

	fields []string // select <fields>
	table  string
}

func (s *SelectCore) From(tableName string) *Select {
	s.table = tableName

	return &Select{child{parent: s}}
}

type queryBuilder struct {
	cols   []identity.Identity
	table  identity.Identity
	wheres []*Where //todo: [][]Where or []OrWhere -> OrWhere{[]AndWhere}

	err error

	indentBuilder *identity.Builder
	strBuilder    *strings.Builder
	tag           string
}

func (qb *queryBuilder) collect(p any) {
	par, ok := p.(parenter)
	if ok {
		qb.collect(par.Parent())
	}

	if par == nil {
		return
	}

	for {
		switch q := p.(type) {
		case *Select:
			return
		case *SelectCore:
			qb.indentBuilder = q.indentBuilder

			if len(q.fields) > 0 {
				qb.cols = make([]identity.Identity, len(q.fields))
				for i := range q.fields {
					qb.cols[i] = qb.indentBuilder.Ident(q.fields[i])
				}
			}

			qb.table = qb.indentBuilder.Ident(q.table)

			return
		case *Where:
			qb.wheres = append(qb.wheres, q)
			return
		default:
			panic(fmt.Sprintf("wrong type in collect %T", p))
		}
	}
}

func (s *Select) ToSQL() (sql string, err error) {
	qb := queryBuilder{
		strBuilder: new(strings.Builder),
	}
	qb.collect(s)

	qb.buildSelectFrom()
	qb.buildWhere()

	return qb.strBuilder.String(), qb.err
}

func (s *Select) ToSQLWithStmts() (sql string, args []any, err error) {
	qb := queryBuilder{
		strBuilder: new(strings.Builder),
	}
	qb.collect(s)
	qb.buildSelectFrom()
	args = qb.buildWherePrepStmt(args)

	return qb.strBuilder.String(), args, qb.err
}

//func (s Select) join(jt joinType, tableName string) joinPart {
//	return joinPart{
//		j: Join{
//			t:         jt,
//			joinTable: s.ident(tableName),
//		},
//		s: s,
//	}
//}
//
//func (s Select) Join(tableName string) joinPart { //nolint:revive
//	return s.join(inner, tableName)
//}
//
//func (s Select) RightJoin(tableName string) joinPart { //nolint:revive
//	return s.join(right, tableName)
//}
//
//func (s Select) LeftJoin(tableName string) joinPart { //nolint:revive
//	return s.join(left, tableName)
//}
//
//func (s Select) FullJoin(tableName string) joinPart { //nolint:revive
//	return s.join(full, tableName)
//}

//func (s Select) Where(columnName string) wherePart[*Select] { //nolint:revive
//	return wherePart[*Select]{
//		field: s.indentBuilder.Ident(columnName),
//		b:      &s,
//	}
//}
//
//func (s Select) Offset(n uint) Select {
//	s.offset = pointer.To(n)
//	return s
//}
//
//func (s Select) Limit(n uint) Select {
//	s.limit = pointer.To(n)
//	return s
//}
//
//func (s Select) OrderBy(fieldName string) orderPart { //nolint:revive
//	return orderPart{
//		field: s.indentBuilder.Ident(fieldName),
//		s:      s,
//	}
//}
//
//func (s *Select) addFieldsFromModel(model any) {
//	rt := reflect.TypeOf(model).Elem()
//	if rt.Kind() != reflect.Struct {
//		s.err = fmt.Errorf("incorrect type of model %s: %w", rt.Kind().String(), ErrValidation)
//		return
//	}
//	for i := 0; i < rt.NumField(); i++ {
//		f, ok := rt.Field(i).Tag.Lookup(s.tag)
//		if !ok {
//			f = stringer2.SnakeCase(rt.Field(i).Name)
//		}
//		s.fields = append(s.fields, s.indentBuilder.Ident(f))
//	}
//}

func (qb *queryBuilder) getFields() string {
	if len(qb.cols) == 0 {
		return "*"
	}

	return strings.Join(qb.cols, ", ")
}

//	func (s *Select) addJoin(j *Join) {
//		s.joins = append(s.joins, j)
//	}
func (qb *queryBuilder) buildSelectFrom() {
	if qb.err != nil {
		return
	}

	_, qb.err = fmt.Fprintf(qb.strBuilder, "SELECT %s FROM %s", qb.getFields(), qb.table)
	//for _, j := range qb.joins {
	//	if qb.err != nil {
	//		return
	//	}
	//	_, qb.err = fmt.Fprint(qb.strBuilder, j.String())
	//}
}

func (qb *queryBuilder) buildWhere() {
	if len(qb.wheres) == 0 {
		return
	}
	strs := make([]string, len(qb.wheres))
	for i := range qb.wheres {
		strs[i] = qb.wheres[i].String(qb.indentBuilder)
	}

	_, qb.err = fmt.Fprintf(qb.strBuilder, " WHERE %s", strings.Join(strs, " AND ")) // todo: build AND and OR separately
}

func (qb *queryBuilder) buildWherePrepStmt(args []any) []any {
	if len(qb.wheres) == 0 {
		return args
	}
	if qb.err != nil {
		return nil
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
