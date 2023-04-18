package query

import (
	"fmt"

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
	}
	panic(c)
}

type Where struct {
	field indent.Indent
	value indent.Value
	cond  Condition
}

func (w Where) String() string {
	return fmt.Sprintf("%s %s %s", w.field.String(), w.cond.String(), w.value.String()) //todo: sanitize value
}

func (w Where) PrepStmtString(num int) (string, any) {
	return fmt.Sprintf("%s %s $%d", w.field.String(), w.cond.String(), num), w.value.Value
}

func (wp wherePart) Equal(v any) Builder {
	wp.b.wheres = append(wp.b.wheres, Where{field: wp.column, value: wp.b.indentBuilder.Value(v), cond: eq})
	return wp.b
}
