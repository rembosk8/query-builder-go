package main

import (
	"fmt"
	"log"

	"github.com/rembosk8/query-builder-go/query/pg"
)

func del() {
	qb := pg.NewQueryBuilder()
	prepQuery := qb.DeleteFrom("example_table").
		Where("year").NotEqual(1990).
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
