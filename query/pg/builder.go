package pg

import (
	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/indent"
	"github.com/rembosk8/query-builder-go/query/pg/sanitize"
)

func NewQueryBuilder() query.Builder {
	indentBuilder := indent.NewBuilder(
		indent.WithIndentSerializer(&sanitize.Indent{}),
		indent.WithValueSerializer(&sanitize.Value{}),
	)

	return query.New(indentBuilder)
}
