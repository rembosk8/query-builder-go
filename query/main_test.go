package query_test

import (
	"os"
	"testing"

	"github.com/rembosk8/query-builder-go/builder/pg"
	"github.com/rembosk8/query-builder-go/query"
)

var qb query.BaseBuilder

func TestMain(m *testing.M) {
	qb = pg.NewQueryBuilder()
	os.Exit(m.Run())
}
