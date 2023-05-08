package main

import (
	"fmt"
	"log"

	"github.com/rembosk8/query-builder-go/query/pg"
)

func main() {
	fmt.Println("Select example:")
	sel()
	fmt.Println("Update example:")
	update()
	fmt.Println("Delete example:")
	del()
	fmt.Println("Insert example:")
	insert()
}

func sel() {
	qb := pg.NewQueryBuilder()
	// generate select query
	prepQuery := qb.Select("id", "name", "year").
		From("example_table").
		Where("year").Between(1990, 2023).
		Where("name").Equal("Max").
		Limit(100).
		Offset(200).
		OrderBy("year").Desc()

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
