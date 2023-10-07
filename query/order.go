package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/internal/identity"
)

type OrderType uint8

const (
	DESC OrderType = iota + 1
	ASC
)

type Order struct {
	child
	field string
	order OrderType
}

func (o *Order) String(idb *identity.Builder) string {
	switch o.order {
	case DESC:
		return fmt.Sprintf("%s DESC", idb.Ident(o.field))
	case ASC:
		return fmt.Sprintf("%s ASC", idb.Ident(o.field))
	default:
		panic("invalid order type")
	}
}

func (o *Order) Desc() *Select {
	o.order = DESC

	return &Select{child{parent: o}}
}

func (o *Order) Asc() *Select {
	o.order = ASC

	return &Select{child{parent: o}}
}
