package benchmark

import (
	"testing"
	"time"

	"github.com/rembosk8/query-builder-go/builder/pg"
	"github.com/rembosk8/query-builder-go/query"
)

func BenchmarkPGBuilder(b *testing.B) {
	qb := pg.NewQueryBuilder()

	getPrepBuild := func() *query.Select {
		return qb.Select("one", "two", "three").
			From("table 1").
			Where("id").Equal(1).
			Where("name").In("n1", "n2", "n3").
			Where("count").Between(1, 100).
			Limit(100).Offset(100)
	}
	preparedQuery := getPrepBuild()

	b.Run("prepare query", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getPrepBuild()
		}
	})

	b.Run("build plain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = preparedQuery.ToSQL()
		}
	})

	b.Run("build with statements", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = preparedQuery.ToSQLWithStmts()
		}
	})
}

type testModel struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func BenchmarkReflectSelect(b *testing.B) {
	qb := pg.NewQueryBuilder()
	m := testModel{}

	prepFunc := func(model any) *query.Select {
		return qb.SelectV2(model).
			From("table 1").
			Where("id").Equal(1).
			Where("name").In("n1", "n2", "n3").
			Where("count").Between(1, 100).
			Limit(100).Offset(100)
	}

	var bdr *query.Select
	for i := 0; i < b.N; i++ {
		bdr = prepFunc(&m)
	}

	_, _, _ = bdr.ToSQLWithStmts() //nolint:dogsled
}
