package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/x/internal/identity"
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
	//op.s.orderBys = append(op.s.orderBys, Order{
	//	field: op.field,
	//	order: DESC,
	//})

	return op.s
}

func (op orderPart) Asc() Select {
	//op.s.orderBys = append(op.s.orderBys, Order{
	//	field: op.field,
	//	order: ASC,
	//})

	return op.s
}

type Order struct {
	field identity.Identity
	order OrderType
}

func (o Order) String(idb *identity.Builder) string {
	switch o.order {
	case DESC:
		return fmt.Sprintf("%s DESC", idb.Ident(o.field))
	case ASC:
		return fmt.Sprintf("%s ASC", idb.Ident(o.field))
	default:
		panic("invalid order type")
	}
}
