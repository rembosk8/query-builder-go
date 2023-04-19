package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/query/indent"
)

type OrderType uint8

const (
	DESC OrderType = iota
	ASC
)

type orderPart struct {
	column indent.Indent
	b      Builder
}

func (op orderPart) Desc() Builder {
	op.b.orderBys = append(op.b.orderBys, Order{
		field: op.column,
		order: DESC,
	})

	return op.b
}

func (op orderPart) Asc() Builder {
	op.b.orderBys = append(op.b.orderBys, Order{
		field: op.column,
		order: ASC,
	})

	return op.b
}

type Order struct {
	field indent.Indent
	order OrderType
}

func (o Order) String() string {
	switch o.order {
	case DESC:
		return fmt.Sprintf("%s DESC", o.field.String())
	case ASC:
		return fmt.Sprintf("%s ASC", o.field.String())
	default:
		panic("invalid order type")
	}
}

var _ fmt.Stringer = &Order{}
