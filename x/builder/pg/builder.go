package pg

import (
	"github.com/rembosk8/query-builder-go/x/internal/identity"
	"github.com/rembosk8/query-builder-go/x/internal/query"
	"github.com/rembosk8/query-builder-go/x/internal/sanitizers/pg"
)

func NewQueryBuilder() query.BaseBuilder {
	indentBuilder := IndentBuilder()

	return query.New(query.WithIdentityBuilder(indentBuilder))
}

func IndentBuilder() *identity.Builder { //todo: check if it should be exported
	return identity.NewBuilder(
		identity.WithIndentSerializer(&pg.Indent{}),
		identity.WithValueSerializer(&pg.Value{}),
	)
}
