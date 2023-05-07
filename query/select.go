package query

import (
	"fmt"
	"reflect"

	"github.com/rembosk8/query-builder-go/helpers/pointer"
	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/indent"
)

var _ sqler = &Select{}

type Select struct {
	baseQuery

	fields   []indent.Indent // select <fields>
	offset   *uint
	limit    *uint
	orderBys []Order
}

func (s Select) ToSql() (sql string, err error) {
	if err = s.initBuild(); err != nil {
		return "", err
	}
	s.buildSqlPlain()

	return s.strBuilder.String(), nil
}

func (s Select) ToSqlWithStmts() (sql string, args []any, err error) {
	if err = s.initBuild(); err != nil {
		return "", nil, err
	}
	args = s.buildPrepStatement()

	return s.strBuilder.String(), args, nil
}

func (s Select) From(tableName string) Select {
	s.setTable(tableName)
	return s
}

func (s Select) Where(columnName string) wherePart[*Select] {
	return wherePart[*Select]{
		column: s.indentBuilder.Indent(columnName),
		b:      &s,
	}
}

func (s Select) Offset(n uint) Select {
	s.offset = pointer.To(n)
	return s
}

func (s Select) Limit(n uint) Select {
	s.limit = pointer.To(n)
	return s
}

func (s Select) OrderBy(fieldName string) orderPart {
	return orderPart{
		column: s.indentBuilder.Indent(fieldName),
		s:      s,
	}
}

func (s *Select) addFieldsFromModel(model any) {
	rt := reflect.TypeOf(model).Elem()
	if rt.Kind() != reflect.Struct {
		s.err = fmt.Errorf("incorrect type of model %s: %w", rt.Kind().String(), ErrValidation)
		return
	}
	for i := 0; i < rt.NumField(); i++ {
		f, ok := rt.Field(i).Tag.Lookup(s.tag)
		if !ok {
			f = stringer.SnakeCase(rt.Field(i).Name)
		}
		s.fields = append(s.fields, s.indentBuilder.Indent(f))
	}
}

func (s *Select) getFields() string {
	if len(s.fields) == 0 {
		return all
	}

	return stringer.Join(s.fields, ", ")
}

func (s *Select) buildSelectFrom() {
	if s.err != nil {
		return
	}

	_, s.err = fmt.Fprintf(s.strBuilder, "SELECT %s FROM %s", s.getFields(), s.table.String())
}

func (s *Select) buildOffset() {
	if s.err != nil {
		return
	}
	if s.offset == nil {
		return
	}

	_, s.err = fmt.Fprintf(s.strBuilder, " OFFSET %d", *s.offset)
}

func (s *Select) buildLimit() {
	if s.err != nil {
		return
	}
	if s.limit == nil {
		return
	}
	_, s.err = fmt.Fprintf(s.strBuilder, " LIMIT %d", *s.limit)
}

func (s *Select) buildOrderBy() {
	if s.err != nil {
		return
	}
	if len(s.orderBys) == 0 {
		return
	}
	_, s.err = fmt.Fprintf(s.strBuilder, " ORDER BY %s", stringer.Join(s.orderBys, ", "))
}

func (s *Select) buildSqlPlain() {
	s.buildSelectFrom()

	s.buildWhere()

	s.buildOrderBy()
	s.buildOffset()
	s.buildLimit()
}

func (s *Select) buildPrepStatement() (args []any) {
	s.buildSelectFrom()

	args = s.buildWherePrepStmt(args)

	s.buildOrderBy()
	s.buildOffset()
	s.buildLimit()

	return args
}
