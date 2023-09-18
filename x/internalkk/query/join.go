package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/x/internalkk/helpers/pointer"
	"github.com/rembosk8/query-builder-go/x/internalkk/identity"
)

type joinType string

const (
	inner = joinType("")
	left  = joinType("LEFT ")
	right = joinType("RIGHT ")
	full  = joinType("FULL ")
)

type joinOn struct {
	t1   string
	col1 string
	t2   string
	col2 string
}

type Join struct {
	child

	t         joinType
	joinTable string
	on        *joinOn // for JOIN <table> ON t1.c1 = t2.c2
	using     *string // for USING (<using>)
}

func (j *Join) String(idb *identity.Builder) string {
	if j.on != nil {
		return fmt.Sprintf(
			" %sJOIN %s ON %s.%s = %s.%s",
			j.t,
			idb.Ident(j.joinTable),
			idb.Ident(j.on.t1),
			idb.Ident(j.on.col1),
			idb.Ident(j.on.t2),
			idb.Ident(j.on.col2),
		)
	}
	if j.using != nil {
		return fmt.Sprintf(
			" %sJOIN %s USING (%s)",
			j.t,
			idb.Ident(j.joinTable),
			idb.Ident(*j.using),
		)
	}
	panic("Join.on nor Join.using is initialized")
}

type joinPart struct {
	j *Join
	s *Select
}

func (jp joinPart) On(table1, column1, table2, column2 string) *Select {
	jp.j.on = &joinOn{
		t1:   table1,
		col1: column1,
		t2:   table2,
		col2: column2,
	}

	return jp.s
}

func (jp joinPart) Using(column string) *Select {
	jp.j.using = pointer.To(column)

	return jp.s
}
