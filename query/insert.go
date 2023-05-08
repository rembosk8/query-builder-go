package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rembosk8/query-builder-go/helpers/stringer"
	"github.com/rembosk8/query-builder-go/query/indent"
)

type Insert struct {
	baseQuery

	fields []indent.Indent
	values []indent.Value
}

var _ sqler = &Insert{}

func (i Insert) ToSql() (sql string, err error) {
	if err = i.initBuild(); err != nil {
		return "", err
	}
	i.buildSqlPlain()

	return i.strBuilder.String(), nil
}

func (i Insert) ToSqlWithStmts() (sql string, args []any, err error) {
	if err = i.initBuild(); err != nil {
		return "", nil, err
	}
	args = i.buildPrepStatement()

	return i.strBuilder.String(), args, nil
}

func (i Insert) Set(field string, value any) Insert {
	i.fields = append(i.fields, i.field(field))
	i.values = append(i.values, i.value(value))

	return i
}

func (i *Insert) buildSqlPlain() {
	i.buildInsertInto()
	i.buildValues()
}

func (i *Insert) buildPrepStatement() (args []any) {
	i.buildInsertInto()
	args = i.buildValueStmts()
	return
}

func (i *Insert) buildInsertInto() {
	if i.err != nil {
		return
	}

	_, i.err = fmt.Fprint(i.strBuilder, "INSERT INTO ", i.table.String())
}

func (i *Insert) buildValues() {
	if i.err != nil {
		return
	}

	_, i.err = fmt.Fprintf(
		i.strBuilder,
		" (%s) VALUES (%s)",
		stringer.Join(i.fields, ", "),
		stringer.Join(i.values, ", "),
	)
}

func (i *Insert) buildValueStmts() (args []any) {
	if i.err != nil {
		return
	}

	numSlice := make([]string, len(i.values))
	args = make([]any, len(i.values))
	for j, v := range i.values {
		numSlice[j] = "$" + strconv.Itoa(j+1)
		args[j] = v.Value
	}

	_, i.err = fmt.Fprintf(
		i.strBuilder,
		" (%s) VALUES (%s)",
		stringer.Join(i.fields, ", "),
		strings.Join(numSlice, ", "),
	)

	return args
}
