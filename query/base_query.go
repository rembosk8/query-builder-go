package query

import (
	"fmt"
	"strings"

	"github.com/rembosk8/query-builder-go/helpers/pointer"
	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/indent"
)

type sqler interface {
	ToSql() (sql string, err error)
	ToSqlWithStmts() (sql string, args []any, err error)

	buildSqlPlain()
	buildPrepStatement() (args []any)
}

type baseQuery struct {
	table  *indent.Indent // from <table>
	wheres []*Where

	err error

	indentBuilder *indent.Builder
	strBuilder    *strings.Builder
	tag           string
}

func (bq *baseQuery) setTable(tblName string) {
	bq.table = pointer.To(bq.indentBuilder.Indent(tblName))
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

func (bq *baseQuery) whereAdd(w *Where) {
	bq.wheres = append(bq.wheres, w)
}

func (bq *baseQuery) value(v any) indent.Value {
	return bq.indentBuilder.Value(v)
}

func (bq *baseQuery) field(f string) indent.Indent {
	return bq.indentBuilder.Indent(f)
}

func (bq *baseQuery) buildWhere() {
	if bq.err != nil {
		return
	}
	if len(bq.wheres) == 0 {
		return
	}
	_, bq.err = fmt.Fprintf(bq.strBuilder, " WHERE %s", stringer.Join(bq.wheres, " AND ")) //todo: build AND and OR separately
}

func (bq *baseQuery) buildWherePrepStmt(args []any) []any {
	if len(bq.wheres) == 0 {
		return args
	}
	if bq.err != nil {
		return nil
	}
	var vals []any
	cnt := len(args) + 1

	_, bq.err = fmt.Fprint(bq.strBuilder, " WHERE ")
	vals, bq.err = bq.wheres[0].PrepStmtString(cnt, bq.strBuilder)
	if bq.err != nil {
		return nil
	}
	args = append(args, vals...)

	cnt += len(vals)
	for i := 1; i < len(bq.wheres); i++ {
		_, bq.err = fmt.Fprint(bq.strBuilder, " AND ")
		if bq.err != nil {
			return nil
		}
		vals, bq.err = bq.wheres[i].PrepStmtString(cnt, bq.strBuilder)
		if bq.err != nil {
			return nil
		}
		args = append(args, vals...)
		cnt += len(vals)
	}

	return args
}
