package query

import "github.com/rembosk8/query-builder-go/query/indent"

type Update struct {
	baseQuery

	fieldValue map[indent.Indent]indent.Value
	returning  []indent.Indent
}
