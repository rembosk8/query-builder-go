package query

//
//import (
//	"errors"
//	"fmt"
//	"reflect"
//	"strings"
//
//	"github.com/rembosk8/query-builder-go/helpers/pointer"
//	"github.com/rembosk8/query-builder-go/helpers/stringer"
//	"github.com/rembosk8/query-builder-go/query/indent"
//)
//
//type Builder struct {
//	fields   []indent.Indent // select <fields>
//	table    *indent.Indent  // from <table>
//	wheres   []*Where
//	offset   *uint
//	limit    *uint
//	orderBys []Order
//
//	err error
//
//	indentBuilder *indent.Builder
//	strBuilder    *strings.Builder
//	tag           string
//}
//
//type BuilderOption func(b *Builder)
//
//func WithStructAnnotationTag(tag string) BuilderOption {
//	return func(b *Builder) {
//		b.tag = tag
//	}
//}
//
//func WithIndentBuilder(ib *indent.Builder) BuilderOption {
//	return func(b *Builder) {
//		b.indentBuilder = ib
//	}
//}
//
//func New(opts ...BuilderOption) Builder {
//	b := Builder{indentBuilder: indent.NewBuilder(), tag: defaultTag}
//
//	for _, o := range opts {
//		o(&b)
//	}
//
//	return b
//}
//
//func (b Builder) Build() (sql string, args []any, err error) {
//	if b.err != nil {
//		return "", nil, err
//	}
//	if b.table == nil {
//		return "", nil, ErrTableNotSet
//	}
//	b.initStringBuilder()
//	args = b.buildPrepStatement()
//
//	return b.strBuilder.String(), args, nil
//}
//
//func (b Builder) BuildPlain() (sql string, err error) {
//	if b.err != nil {
//		return "", err
//	}
//	if b.table == nil {
//		return "", ErrTableNotSet
//	}
//	b.initStringBuilder()
//	b.buildSqlPlain()
//	if b.err != nil {
//		return "", err
//	}
//
//	return b.strBuilder.String(), nil
//}
//
//func (b Builder) From(tableName string) Builder {
//	b.table = pointer.To(b.indentBuilder.Indent(tableName))
//	return b
//}
//
//func (b Builder) Select(fields ...string) Builder {
//	for _, f := range fields {
//		b.fields = append(b.fields, b.indentBuilder.Indent(f))
//	}
//	return b
//}
//
//func (b *Builder) addFieldsFromModel(model any) {
//	rt := reflect.TypeOf(model).Elem()
//	if rt.Kind() != reflect.Struct {
//		b.err = fmt.Errorf("incorrect type of model %s: %w", rt.Kind().String(), ErrValidation)
//		return
//	}
//	for i := 0; i < rt.NumField(); i++ {
//		f, ok := rt.Field(i).Tag.Lookup(b.tag)
//		if !ok {
//			f = stringer.SnakeCase(rt.Field(i).Name)
//		}
//		b.fields = append(b.fields, b.indentBuilder.Indent(f))
//	}
//}
//
//func (b Builder) SelectV2(model any) Builder {
//	if b.err != nil {
//		return b
//	}
//	b.addFieldsFromModel(model)
//
//	return b
//}
//
//func (b Builder) Where(columnName string) wherePart[*Builder] {
//	return wherePart[*Builder]{
//		column: b.indentBuilder.Indent(columnName),
//		b:      &b,
//	}
//}
//
//func (b Builder) Offset(n uint) Builder {
//	b.offset = pointer.To(n)
//	return b
//}
//
//func (b Builder) Limit(n uint) Builder {
//	b.limit = pointer.To(n)
//	return b
//}
//
//func (b Builder) OrderBy(fieldName string) orderPart {
//	return orderPart{
//		column: b.indentBuilder.Indent(fieldName),
//		s:      b,
//	}
//}
//
//func (b *Builder) initStringBuilder() {
//	b.strBuilder = new(strings.Builder)
//}
//
//func (b *Builder) getFields() string {
//	if len(b.fields) == 0 {
//		return all
//	}
//
//	return stringer.Join(b.fields, ", ")
//}
//
//func (b *Builder) buildSelectFrom() {
//	if b.err != nil {
//		return
//	}
//
//	_, b.err = fmt.Fprintf(b.strBuilder, "SELECT %s FROM %s", b.getFields(), b.table.String())
//}
//
//func (b *Builder) buildWherePlain() {
//	if b.err != nil {
//		return
//	}
//	if len(b.wheres) == 0 {
//		return
//	}
//	_, b.err = fmt.Fprintf(b.strBuilder, " WHERE %s", stringer.Join(b.wheres, " AND ")) //todo: build AND and OR separately
//}
//
//func (b *Builder) buildOffset() {
//	if b.err != nil {
//		return
//	}
//	if b.offset == nil {
//		return
//	}
//
//	_, b.err = fmt.Fprintf(b.strBuilder, " OFFSET %d", *b.offset)
//}
//
//func (b *Builder) buildLimit() {
//	if b.err != nil {
//		return
//	}
//	if b.limit == nil {
//		return
//	}
//	_, b.err = fmt.Fprintf(b.strBuilder, " LIMIT %d", *b.limit)
//}
//
//func (b *Builder) buildWherePrepStmt() []any {
//	if len(b.wheres) == 0 {
//		return nil
//	}
//	if b.err != nil {
//		return nil
//	}
//	var args, vals []any
//	cnt := 1
//
//	_, b.err = fmt.Fprint(b.strBuilder, " WHERE ")
//	vals, b.err = b.wheres[0].PrepStmtString(cnt, b.strBuilder)
//	if b.err != nil {
//		return nil
//	}
//	args = append(args, vals...)
//
//	cnt += len(vals)
//	for i := 1; i < len(b.wheres); i++ {
//		_, b.err = fmt.Fprint(b.strBuilder, " AND ")
//		if b.err != nil {
//			return nil
//		}
//		vals, b.err = b.wheres[i].PrepStmtString(cnt, b.strBuilder)
//		if b.err != nil {
//			return nil
//		}
//		args = append(args, vals...)
//		cnt += len(vals)
//	}
//
//	return args
//}
//
//func (b *Builder) buildSqlPlain() {
//	b.buildSelectFrom()
//
//	b.buildWherePlain()
//
//	b.buildOrderBy()
//	b.buildOffset()
//	b.buildLimit()
//}
//
//func (b *Builder) buildPrepStatement() (args []any) {
//	b.buildSelectFrom()
//
//	args = b.buildWherePrepStmt()
//
//	b.buildOrderBy()
//	b.buildOffset()
//	b.buildLimit()
//
//	return args
//}
//
//func (b *Builder) buildOrderBy() {
//	if b.err != nil {
//		return
//	}
//	if len(b.orderBys) == 0 {
//		return
//	}
//	_, b.err = fmt.Fprintf(b.strBuilder, " ORDER BY %s", stringer.Join(b.orderBys, ", "))
//}
//
//func (b *Builder) whereAdd(w *Where) {
//	b.wheres = append(b.wheres, w)
//}
//
//func (b Builder) value(v any) indent.Value {
//	return b.indentBuilder.Value(v)
//}
