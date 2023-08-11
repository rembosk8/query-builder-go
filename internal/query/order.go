package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/internal/identity"
)

type OrderType uint8

const (
	DESC OrderType = iota
	ASC
)

type orderPart struct {
	column identity.Identity
	s      Select
}

func (op orderPart) Desc() Select {
	op.s.orderBys = append(op.s.orderBys, Order{
		field: op.column,
		order: DESC,
	})

	return op.s
}

func (op orderPart) Asc() Select {
	op.s.orderBys = append(op.s.orderBys, Order{
		field: op.column,
		order: ASC,
	})

	return op.s
}

type Order struct {
	field identity.Identity
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
