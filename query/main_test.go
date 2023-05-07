package query_test

import (
	"os"
	"testing"

	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/pg"
)

var qb query.BaseBuilder

func TestMain(m *testing.M) {
	qb = pg.NewQueryBuilder()
	os.Exit(m.Run())
}
