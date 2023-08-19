package query_test

import (
	"os"
	"testing"

	"github.com/rembosk8/query-builder-go/x/builder/pg"
	"github.com/rembosk8/query-builder-go/x/internal/query"
)

var qb query.BaseBuilder

func TestMain(m *testing.M) {
	qb = pg.NewQueryBuilder()
	os.Exit(m.Run())
}
