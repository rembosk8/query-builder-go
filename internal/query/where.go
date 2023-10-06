package query

import (
	"fmt"
	"io"
	"strings"

	"github.com/rembosk8/query-builder-go/internal/identity"
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

const numsStr = " $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39"

func getNums(s, cnt int) string {
	if s == 0 {
		return ""
	}
	s--
	start := s * 4
	pos := start + (cnt-1)*4 + 3

	if s > 9 {
		start += s - 9
	}
	if s+cnt > 9 {
		pos += s + cnt - 9
	}

	return numsStr[start+1 : pos]
}

func (w *Where) PrepStmtString(num int, wr io.Writer, idBuilder *identity.Builder) ([]any, error) {
	switch w.cond {
	case eq, ne, le, lq, gt, gq, like, notLike:
		if _, err := fmt.Fprintf(wr, "%s %s $%d", idBuilder.Ident(w.field), w.cond.String(), num); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}
		return w.value, nil
	case in, notIn:
		if _, err := fmt.Fprintf(wr, "%s %s (%s)", idBuilder.Ident(w.field), w.cond.String(), getNums(num, len(w.value))); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return w.value, nil
	case isNull, isNotNull:
		if _, err := fmt.Fprintf(wr, "%s %s", idBuilder.Ident(w.field), w.cond.String()); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return nil, nil
	case between, notBetween:
		if _, err := fmt.Fprintf(wr, "%s %s $%d AND $%d", idBuilder.Ident(w.field), w.cond.String(), num, num+1); err != nil {
			return nil, fmt.Errorf("write into writer: %w", err)
		}

		return w.value, nil
	}

	panic("unknown where condition")
}

func (w *WherePart[T]) Equal(v any) T {
	w.Where.cond = eq
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) NotEqual(v any) T {
	w.Where.cond = ne
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) Less(v any) T {
	w.Where.cond = le
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) LessEqual(v any) T {
	w.Where.cond = lq
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) Greater(v any) T {
	w.Where.cond = gt
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) GreaterEqual(v any) T {
	w.Where.cond = gq
	w.Where.value = append(w.Where.value, v)

	return w.b
}

func (w *WherePart[T]) In(vs ...any) T {
	w.Where.cond = in
	w.Where.value = append(w.Where.value, vs...)

	return w.b
}

func (w *WherePart[T]) NotIn(vs ...any) T {
	w.Where.cond = notIn
	w.Where.value = append(w.Where.value, vs...)

	return w.b
}

func (w *WherePart[T]) IsNull() T {
	w.Where.cond = isNull

	return w.b
}

func (w *WherePart[T]) IsNotNull() T {
	w.Where.cond = isNotNull

	return w.b
}

func (w *WherePart[T]) Between(start, end any) T {
	w.Where.value = append(w.Where.value, start, end)
	w.Where.cond = between

	return w.b
}

func (w *WherePart[T]) NotBetween(start, end any) T {
	w.Where.value = append(w.Where.value, start, end)
	w.Where.cond = notBetween

	return w.b
}

func (w *WherePart[T]) Like(pattern string) T {
	w.Where.cond = like
	w.Where.value = append(w.Where.value, pattern)

	return w.b
}

func (w *WherePart[T]) NotLike(pattern string) T {
	w.Where.cond = notLike
	w.Where.value = append(w.Where.value, pattern)

	return w.b
}
