package query

import (
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/indent"
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
	notIn
	isNull
	isNotNull
	between
	notBetween
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
	case notIn:
		return "NOT IN"
	case isNull:
		return "IS NULL"
	case isNotNull:
		return "IS NOT NULL"
	case between:
		return "BETWEEN"
	case notBetween:
		return "NOT BETWEEN"
	}
	panic(c)
}

type Where struct {
	field indent.Indent
	value []indent.Value
	cond  Condition
}

func (w *Where) String() string {
	switch w.cond {
	case eq, ne, le, lq, gt, gq:
		return fmt.Sprintf("%s %s %s", w.field.String(), w.cond.String(), w.value[0].String())
	case in, notIn:
		return fmt.Sprintf("%s %s (%s)", w.field.String(), w.cond.String(), stringer.Join(w.value, ", "))
	case isNull, isNotNull:
		return fmt.Sprintf("%s %s", w.field.String(), w.cond.String())
	case between, notBetween:
		return fmt.Sprintf("%s %s %s AND %s", w.field.String(), w.cond.String(), w.value[0].String(), w.value[1].String())
	}

	panic("unknown where condition")
}

func (w *Where) PrepStmtString(num int) (string, []any) {
	vals := make([]any, len(w.value))
	for i := range w.value {
		vals[i] = w.value[i].Value
	}

	switch w.cond {
	case eq, ne, le, lq, gt, gq:
		return fmt.Sprintf("%s %s $%d", w.field.String(), w.cond.String(), num), vals
	case in, notIn:
		nums := make([]string, len(vals))
		for i := range vals {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}

		return fmt.Sprintf("%s %s (%s)", w.field.String(), w.cond.String(), strings.Join(nums, ", ")), vals
	case isNull, isNotNull:
		return fmt.Sprintf("%s %s", w.field.String(), w.cond.String()), nil
	case between, notBetween:
		nums := make([]string, len(vals))
		for i := range vals {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}
		return fmt.Sprintf("%s %s %s AND %s", w.field.String(), w.cond.String(), nums[0], nums[1]), vals
	}

	panic("unknown where condition")
}

func (wp wherePart) Equal(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: eq})
	return wp.b
}

func (wp wherePart) NotEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: ne})
	return wp.b
}

func (wp wherePart) Less(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: le})
	return wp.b
}

func (wp wherePart) LessEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: lq})
	return wp.b
}

func (wp wherePart) Greater(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: gt})
	return wp.b
}

func (wp wherePart) GreaterEqual(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: []indent.Value{wp.b.indentBuilder.Value(v)}, cond: gq})
	return wp.b
}

func (wp wherePart) In(vs ...any) Builder {
	values := make([]indent.Value, len(vs))
	for i := range vs {
		values[i] = wp.b.indentBuilder.Value(vs[i])
	}
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: values, cond: in})
	return wp.b
}

func (wp wherePart) NotIn(vs ...any) Builder {
	values := make([]indent.Value, len(vs))
	for i := range vs {
		values[i] = wp.b.indentBuilder.Value(vs[i])
	}
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: values, cond: notIn})
	return wp.b
}

func (wp wherePart) IsNull() Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: nil, cond: isNull})
	return wp.b
}

func (wp wherePart) IsNotNull() Builder {
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: nil, cond: isNotNull})
	return wp.b
}

func (wp wherePart) Between(start, end any) Builder {
	values := []indent.Value{
		wp.b.indentBuilder.Value(start),
		wp.b.indentBuilder.Value(end),
	}
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: values, cond: between})
	return wp.b
}

func (wp wherePart) NotBetween(start, end any) Builder {
	values := []indent.Value{
		wp.b.indentBuilder.Value(start),
		wp.b.indentBuilder.Value(end),
	}
	wp.b.wheres = append(wp.b.wheres, &Where{field: wp.column, value: values, cond: notBetween})
	return wp.b
}
