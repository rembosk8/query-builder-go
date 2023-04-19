package query

import (
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/query/indent"
	"github.com/rembosk8/query-builder-go/stringer"
)

type wherePart struct {
	column indent.Indent
	b      Builder
}

type Condition int16

const (
	eq Condition = iota
	ne
	le
	lq
	gt
	gq
	in
)

func (c Condition) String() string {
	switch c {
	case eq:
		return "="
	case ne:
		return "!="
	case le:
		return "<"
	case lq:
		return "<="
	case gt:
		return ">"
	case gq:
		return ">="
	case in:
		return "IN"
	}
	panic(c)
}

type Where struct {
	field indent.Indent
	value []indent.Value
	cond  Condition
}

func (w Where) String() string {
	switch w.cond {
	case eq, ne, le, lq, gt, gq:
		return fmt.Sprintf("%s %s %s", w.field.String(), w.cond.String(), w.value[0].String())
	case in:
		return fmt.Sprintf("%s %s (%s)", w.field.String(), w.cond.String(), stringer.Join(w.value, ", "))
	}

	panic("unknown where condition")
}

func (w Where) PrepStmtString(num int) (string, []any) {
	vals := make([]any, len(w.value))
	for i := range w.value {
		vals[i] = w.value[i].Value
	}

	switch w.cond {
	case eq, ne, le, lq, gt, gq:
		return fmt.Sprintf("%s %s $%d", w.field.String(), w.cond.String(), num), vals
	case in:
		nums := make([]string, len(vals))
		for i := range vals {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}

		return fmt.Sprintf("%s %s (%s)", w.field.String(), w.cond.String(), strings.Join(nums, ", ")), vals
	}

	panic("unknown where condition")
}

func (wp wherePart) Equal(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: eq})
	return wp.b
}

func (wp wherePart) NotEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: ne})
	return wp.b
}

func (wp wherePart) Less(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: le})
	return wp.b
}

func (wp wherePart) LessEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: lq})
	return wp.b
}

func (wp wherePart) Greater(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: gt})
	return wp.b
}

func (wp wherePart) GreaterEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: gq})
	return wp.b
}

func (wp wherePart) In(vs ...any) Builder {
	values := make([]indent.Value, len(vs))
	for i := range vs {
		values[i] = wp.b.indentBuilder.Value(vs[i])
	}
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: values, cond: in})
	return wp.b
}
