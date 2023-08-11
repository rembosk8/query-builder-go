package query

import (
	"fmt"
	"io"
	"strings"

	"github.com/rembosk8/query-builder-go/internal/helpers/stringer"
	"github.com/rembosk8/query-builder-go/internal/identity"
)

type whereAdder interface {
	whereAdd(where *Where)
	value(v any) identity.Value
}

type wherePart[T whereAdder] struct {
	column identity.Identity
	b      T
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
	like
	notLike
)

var conditionStrings = []string{ //nolint:gochecknoglobals
	eq:         "=",
	ne:         "!=",
	le:         "<",
	lq:         "<=",
	gt:         ">",
	gq:         ">=",
	in:         "IN",
	notIn:      "NOT IN",
	isNull:     "IS NULL",
	isNotNull:  "IS NOT NULL",
	between:    "BETWEEN",
	notBetween: "NOT BETWEEN",
	like:       "LIKE",
	notLike:    "NOT LIKE",
}

func (c Condition) String() string {
	return conditionStrings[c]
}

type Where struct {
	field identity.Identity
	value []identity.Value
	cond  Condition
}

func (w *Where) String() string {
	switch w.cond {
	case eq, ne, le, lq, gt, gq, like, notLike:
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

func (w *Where) PrepStmtString(num int, wr io.Writer) ([]any, error) {
	vals := make([]any, len(w.value))
	for i := range w.value {
		vals[i] = w.value[i].Value
	}

	switch w.cond {
	case eq, ne, le, lq, gt, gq, like, notLike:
		if _, err := fmt.Fprintf(wr, "%s %s $%d", w.field.String(), w.cond.String(), num); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}
		return vals, nil
	case in, notIn:
		nums := make([]string, len(vals))
		for i := range vals {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}
		if _, err := fmt.Fprintf(wr, "%s %s (%s)", w.field.String(), w.cond.String(), strings.Join(nums, ", ")); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return vals, nil
	case isNull, isNotNull:
		if _, err := fmt.Fprintf(wr, "%s %s", w.field.String(), w.cond.String()); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return nil, nil
	case between, notBetween:
		nums := make([]string, len(vals))
		for i := range vals {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}

		if _, err := fmt.Fprintf(wr, "%s %s %s AND %s", w.field.String(), w.cond.String(), nums[0], nums[1]); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return vals, nil
	}

	panic("unknown where condition")
}

func (wp wherePart[T]) Equal(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: eq})

	return wp.b
}

func (wp wherePart[T]) NotEqual(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: ne})

	return wp.b
}

func (wp wherePart[T]) Less(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: le})

	return wp.b
}

func (wp wherePart[T]) LessEqual(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: lq})

	return wp.b
}

func (wp wherePart[T]) Greater(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: gt})

	return wp.b
}

func (wp wherePart[T]) GreaterEqual(v any) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(v)}, cond: gq})

	return wp.b
}

func (wp wherePart[T]) In(vs ...any) T {
	values := make([]identity.Value, len(vs))
	for i := range vs {
		values[i] = wp.b.value(vs[i])
	}
	wp.b.whereAdd(&Where{field: wp.column, value: values, cond: in})

	return wp.b
}

func (wp wherePart[T]) NotIn(vs ...any) T {
	values := make([]identity.Value, len(vs))
	for i := range vs {
		values[i] = wp.b.value(vs[i])
	}
	wp.b.whereAdd(&Where{field: wp.column, value: values, cond: notIn})

	return wp.b
}

func (wp wherePart[T]) IsNull() T {
	wp.b.whereAdd(&Where{field: wp.column, value: nil, cond: isNull})

	return wp.b
}

func (wp wherePart[T]) IsNotNull() T {
	wp.b.whereAdd(&Where{field: wp.column, value: nil, cond: isNotNull})

	return wp.b
}

func (wp wherePart[T]) Between(start, end any) T {
	values := []identity.Value{
		wp.b.value(start),
		wp.b.value(end),
	}
	wp.b.whereAdd(&Where{field: wp.column, value: values, cond: between})

	return wp.b
}

func (wp wherePart[T]) NotBetween(start, end any) T {
	values := []identity.Value{
		wp.b.value(start),
		wp.b.value(end),
	}
	wp.b.whereAdd(&Where{field: wp.column, value: values, cond: notBetween})

	return wp.b
}

func (wp wherePart[T]) Like(pattern string) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(pattern)}, cond: like})

	return wp.b
}

func (wp wherePart[T]) NotLike(pattern string) T {
	wp.b.whereAdd(&Where{field: wp.column, value: []identity.Value{wp.b.value(pattern)}, cond: notLike})

	return wp.b
}
