package query

import (
	"strings"

	"github.com/rembosk8/query-builder-go/x/internal/helpers/pointer"
	identity2 "github.com/rembosk8/query-builder-go/x/internal/identity"
)

type Builder interface {
	ToSQL() (sql string, err error)
	ToSQLWithStmts() (sql string, args []any, err error)
}

type baseQuery struct {
	table *identity2.Identity // from <table>

	err error

	indentBuilder *identity2.Builder
	strBuilder    *strings.Builder
	tag           string
}

func (bq *baseQuery) setTable(tblName string) {
	bq.table = pointer.To(bq.indentBuilder.Ident(tblName))
}

func (bq *baseQuery) initBuild() error {
	if bq.err != nil {
		return bq.err
	}
	if bq.table == nil {
		return ErrTableNotSet
	}
	bq.strBuilder = new(strings.Builder)

	return nil
}

func (bq *baseQuery) value(v any) identity2.Value {
	return bq.indentBuilder.Value(v)
}

func (bq *baseQuery) ident(f string) identity2.Identity {
	return bq.indentBuilder.Ident(f)
}

func (bq *baseQuery) buildWhere() {
	if bq.err != nil {
		return
	}
	//if len(bq.wheres) == 0 {
	//	return
	//}
	//_, bq.err = fmt.Fprintf(bq.strBuilder, " WHERE %s", stringer.Join(bq.wheres, " AND ")) // todo: build AND and OR separately
}

func (bq *baseQuery) buildWherePrepStmt(args []any) []any {
	//if len(bq.wheres) == 0 {
	//	return args
	//}
	//if bq.err != nil {
	//	return nil
	//}
	//var vals []any
	//cnt := len(args) + 1
	//
	//_, bq.err = fmt.Fprint(bq.strBuilder, " WHERE ")
	//vals, bq.err = bq.wheres[0].PrepStmtString(cnt, bq.strBuilder)
	//if bq.err != nil {
	//	return nil
	//}
	//args = append(args, vals...)
	//
	//cnt += len(vals)
	//for i := 1; i < len(bq.wheres); i++ {
	//	_, bq.err = fmt.Fprint(bq.strBuilder, " AND ")
	//	if bq.err != nil {
	//		return nil
	//	}
	//	vals, bq.err = bq.wheres[i].PrepStmtString(cnt, bq.strBuilder)
	//	if bq.err != nil {
	//		return nil
	//	}
	//	args = append(args, vals...)
	//	cnt += len(vals)
	//}

	return args
}
