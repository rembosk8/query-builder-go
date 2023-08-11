package query

import (
	"fmt"

	"github.com/rembosk8/query-builder-go/internal/helpers/pointer"
	"github.com/rembosk8/query-builder-go/internal/identity"
)

type joinType string

const (
	inner = joinType("")
	left  = joinType("LEFT ")
	right = joinType("RIGHT ")
	full  = joinType("FULL ")
)

type joinOn struct {
	t1   identity.Identity
	col1 identity.Identity
	t2   identity.Identity
	col2 identity.Identity
}

type Join struct {
	t         joinType
	joinTable identity.Identity
	on        *joinOn            // for JOIN <table> ON t1.c1 = t2.c2
	using     *identity.Identity // for USING (<using>)
}

func (j *Join) String() string {
	if j.on != nil {
		return fmt.Sprintf(
			" %sJOIN %s ON %s.%s = %s.%s",
			j.t,
			j.joinTable.String(),
			j.on.t1,
			j.on.col1,
			j.on.t2,
			j.on.col2,
		)
	}
	if j.using != nil {
		return fmt.Sprintf(
			" %sJOIN %s USING (%s)",
			j.t,
			j.joinTable.String(),
			j.using.String(),
		)
	}
	panic("Join.on nor Join.using is initialized")
}

var _ fmt.Stringer = &Join{}

type joinPart struct {
	j Join
	s Select
}

func (jp joinPart) On(table1, column1, table2, column2 string) Select {
	jp.j.on = &joinOn{
		t1:   jp.s.ident(table1),
		col1: jp.s.ident(column1),
		t2:   jp.s.ident(table2),
		col2: jp.s.ident(column2),
	}
	jp.s.addJoin(&jp.j)

	return jp.s
}

func (jp joinPart) Using(column string) Select {
	jp.j.using = pointer.To(jp.s.ident(column))
	jp.s.addJoin(&jp.j)

	return jp.s
}
