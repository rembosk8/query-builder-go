package query

import (
	"fmt"
	"io"
	"strings"

	"github.com/rembosk8/query-builder-go/x/internal/identity"
)

type WherePart[T Builder] struct {
	*Where
	b T
}

type Where struct {
	child
	field string
	cond  condition
	value []any
}

type condition int16

const (
	eq condition = iota
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

func (c condition) String() string {
	return conditionStrings[c]
}

func (w *Where) String(idBuilder *identity.Builder) string {
	switch w.cond {
	case eq, ne, le, lq, gt, gq, like, notLike:
		return fmt.Sprintf("%s %s %s", idBuilder.Ident(w.field), w.cond.String(), idBuilder.Value(w.value[0]))
	case in, notIn:
		return fmt.Sprintf("%s %s (%s)", idBuilder.Ident(w.field), w.cond.String(), strings.Join(idBuilder.Values(w.value), ", "))
	case isNull, isNotNull:
		return fmt.Sprintf("%s %s", idBuilder.Ident(w.field), w.cond.String())
	case between, notBetween:
		return fmt.Sprintf("%s %s %s AND %s", idBuilder.Ident(w.field), w.cond.String(), idBuilder.Value(w.value[0]), idBuilder.Value(w.value[1]))
	}

	panic("unknown where condition")
}

func (w *Where) PrepStmtString(num int, wr io.Writer, idBuilder *identity.Builder) ([]any, error) {
	switch w.cond {
	case eq, ne, le, lq, gt, gq, like, notLike:
		if _, err := fmt.Fprintf(wr, "%s %s $%d", idBuilder.Ident(w.field), w.cond.String(), num); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}
		return w.value, nil
	case in, notIn:
		nums := make([]string, len(w.value))
		for i := range w.value {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}
		if _, err := fmt.Fprintf(wr, "%s %s (%s)", idBuilder.Ident(w.field), w.cond.String(), strings.Join(nums, ", ")); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return w.value, nil
	case isNull, isNotNull:
		if _, err := fmt.Fprintf(wr, "%s %s", idBuilder.Ident(w.field), w.cond.String()); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return nil, nil
	case between, notBetween:
		nums := make([]string, len(w.value))
		for i := range w.value {
			nums[i] = fmt.Sprintf("$%d", num)
			num++
		}

		if _, err := fmt.Fprintf(wr, "%s %s %s AND %s", idBuilder.Ident(w.field), w.cond.String(), nums[0], nums[1]); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return w.value, nil
	}

	panic("unknown where condition")
}

func (w *WherePart[T]) Equal(v any) T {
	w.Where.cond = eq
	w.Where.value = []any{v}

	return w.b
}

func (w *WherePart[T]) NotEqual(v any) T {
	w.cond = ne
	w.value = []any{v}

	return w.b
}

func (w *WherePart[T]) Less(v any) T {
	w.cond = le
	w.value = []any{v}

	return w.b
}

func (w *WherePart[T]) LessEqual(v any) T {
	w.cond = lq
	w.value = []any{v}

	return w.b
}

func (w *WherePart[T]) Greater(v any) T {
	w.cond = gt
	w.value = []any{v}

	return w.b
}

func (w *WherePart[T]) GreaterEqual(v any) T {
	w.cond = gq
	w.value = []any{v}

	return w.b
}

func (w *WherePart[T]) In(vs ...any) T {
	w.cond = in
	w.value = vs

	return w.b
}

func (w *WherePart[T]) NotIn(vs ...any) T {
	w.cond = notIn
	w.value = vs

	return w.b
}

func (w *WherePart[T]) IsNull() T {
	w.cond = isNull

	return w.b
}

func (w *WherePart[T]) IsNotNull() T {
	w.cond = isNotNull

	return w.b
}

func (w *WherePart[T]) Between(start, end any) T {
	w.value = []any{
		start,
		end,
	}
	w.cond = between

	return w.b
}

func (w *WherePart[T]) NotBetween(start, end any) T {
	w.value = []any{
		start,
		end,
	}
	w.cond = notBetween

	return w.b
}

func (w *WherePart[T]) Like(pattern string) T {
	w.cond = like
	w.value = []any{pattern}

	return w.b
}

func (w *WherePart[T]) NotLike(pattern string) T {
	w.cond = notLike
	w.value = []any{pattern}

	return w.b
}
