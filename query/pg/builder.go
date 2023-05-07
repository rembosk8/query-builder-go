package pg

import (
	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/indent"
	"github.com/rembosk8/query-builder-go/query/pg/sanitize"
)

func NewQueryBuilder() query.BaseBuilder {
	indentBuilder := IndentBuilder()

	return query.New(query.WithIndentBuilder(indentBuilder))
}

func IndentBuilder() *indent.Builder {
	return indent.NewBuilder(
		indent.WithIndentSerializer(&sanitize.Indent{}),
		indent.WithValueSerializer(&sanitize.Value{}),
	)
}
