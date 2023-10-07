package main_test

import (
	"log"
	"testing"

	"github.com/rembosk8/query-builder-go/builder/pg"
)

func TestA(t *testing.T) {
	for i := 0; i < 1000; i++ {
		qb := pg.NewQueryBuilder()
		prepQuery := qb.Select("id", "name", "year").
			From("example_table").
			Where("year").Between(1990, 2023).
			Where("name").Equal("Max").
			Limit(100).
			Offset(200).
			OrderBy("year").Desc()
		sql, err := prepQuery.ToSQL()
		if err != nil {
			log.Fatal(err)
		}
		print(sql)
	}

}
