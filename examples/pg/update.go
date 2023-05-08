package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rembosk8/query-builder-go/query/pg"
)

func update() {
	qb := pg.NewQueryBuilder()
	prepQuery := qb.Update("example_table").
		Set("name", "John").
		Set("updated_at", time.Now()).
		Where("name").Equal("Jjjohn").
		Returning("id")

	sql, err := prepQuery.ToSql()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL plain request: ")
	fmt.Println(sql)

	sql, args, err := prepQuery.ToSqlWithStmts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL with placeholders request: ")
	fmt.Println(sql)

	fmt.Println("List of arguments for placeholders: ")
	fmt.Println(args...)
}
