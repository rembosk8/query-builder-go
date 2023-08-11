package main

import (
	"fmt"
	"log"

	"github.com/rembosk8/query-builder-go/builder/pg"
)

func insert() {
	qb := pg.NewQueryBuilder()
	prepQuery := qb.InsertInto("example_table").
		Set("name", "John").
		Set("year", 1989)

	sql, err := prepQuery.ToSQL()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL plain request: ")
	fmt.Println(sql)

	sql, args, err := prepQuery.ToSQLWithStmts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL with placeholders request: ")
	fmt.Println(sql)

	fmt.Println("List of arguments for placeholders: ")
	fmt.Println(args...)
}
