package pg

import (
	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/identity"
	"github.com/rembosk8/query-builder-go/query/pg/sanitize"
)

func NewQueryBuilder() query.BaseBuilder {
	indentBuilder := IndentBuilder()

	return query.New(query.WithIdentityBuilder(indentBuilder))
}

func IndentBuilder() *identity.Builder {
	return identity.NewBuilder(
		identity.WithIndentSerializer(&sanitize.Indent{}),
		identity.WithValueSerializer(&sanitize.Value{}),
	)
}
