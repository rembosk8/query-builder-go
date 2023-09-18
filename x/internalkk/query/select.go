package query

import (
	"strings"

	"github.com/rembosk8/query-builder-go/x/internalkk/identity"
)

type selecter interface {
	Builder

	Where(string) *WherePart[*Select]
}

type parenter interface {
	Parent() any //todo: check with any instead of parenter
}

type child struct {
	parent any
}

func (c child) Parent() any {
	return c.parent
}

type Select struct {
	child
}

func (s Select) Where(field string) *WherePart[*Select] {
	// todo: try pointer receiver and check benchmarks
	// todo: check heap move
	w := Where{
		child: child{parent: s.parent},
		field: field,
	}

	s.parent = &w

	return &WherePart[*Select]{
		b:     &s,
		Where: &w,
	}
}

var _ selecter = &Select{}

type core struct {
	child
	indentBuilder *identity.Builder
}

func qbInit(c any) *queryBuilder { //todo: check how to replace 'any'
	qb := &queryBuilder{
		// todo: try to init identityBuilder from here, not from collect
		strBuilder: &strings.Builder{},
	}
	qb.collect(c)

	return qb
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

func (s Select) ToSQL() (sql string, err error) {
	qb := qbInit(s)

	qb.buildSelectFrom()
	qb.buildWhere()
	qb.BuildOrderBy()
	qb.BuildOffset()
	qb.BuildLimit()

	return qb.strBuilder.String(), qb.err
}

func (s Select) ToSQLWithStmts() (sql string, args []any, err error) {
	qb := qbInit(s)
	qb.buildSelectFrom()
	args = qb.buildWherePrepStmt(args)
	qb.BuildOrderBy()
	qb.BuildOffset()
	qb.BuildLimit()

	return qb.strBuilder.String(), args, qb.err
}

func (s Select) join(jt joinType, tableName string) *joinPart {
	j := Join{
		child:     child{parent: s.parent},
		t:         jt,
		joinTable: tableName,
	}

	jp := &joinPart{
		j: &j,
		s: &s,
	}

	s.parent = &j

	return jp
}

func (s Select) Join(tableName string) *joinPart { //nolint:revive
	return s.join(inner, tableName)
}

func (s Select) RightJoin(tableName string) *joinPart { //nolint:revive
	return s.join(right, tableName)
}

func (s Select) LeftJoin(tableName string) *joinPart { //nolint:revive
	return s.join(left, tableName)
}

func (s Select) FullJoin(tableName string) *joinPart { //nolint:revive
	return s.join(full, tableName)
}

type Offset struct {
	child
	offset uint
}

func (s Select) Offset(n uint) *Select {
	o := &Offset{
		child:  child{parent: s.parent},
		offset: n,
	}
	s.parent = o

	return &s
}

type Limit struct {
	child
	limit uint
}

func (s Select) Limit(n uint) *Select {
	l := &Limit{
		child: child{parent: s.parent},
		limit: n,
	}
	s.parent = l

	return &s
}

func (s Select) OrderBy(fieldName string) *Order { //nolint:revive
	return &Order{
		child: child{parent: s.parent},
		field: fieldName,
	}
}

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
