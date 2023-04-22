package benchmark

import (
	"testing"

	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/pg"
)

var (
	preparedQuery query.Builder
)

func BenchmarkPGBuilder(b *testing.B) {

	qb := pg.NewQueryBuilder()

	getPrepBuild := func() query.Builder {
		return qb.Select("one", "two", "three").
			From("table 1").
			Where("id").Equal(1).
			Where("name").In("n1", "n2", "n3").
			Where("count").Between(1, 100).
			Limit(100).Offset(100)
	}
	preparedQuery = getPrepBuild()

	b.Run("prepare query", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getPrepBuild()
		}
	})

	b.Run("build plain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = preparedQuery.BuildPlain()
		}
	})

	b.Run("build with statements", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = preparedQuery.Build()
		}
	})

}
