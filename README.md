# query-builder-go
SQL Query Bulder for golang which helps to build SQL request dynamically

# Examples

## Select [example.go](examples/pg/select.go)
### Code
```go
package main

import (
	"fmt"
	"log"

	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/pg"
)

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
```
### Output
```
Generated SQL plain request: 
SELECT "id", "name", "year" FROM "example_table" WHERE "year" BETWEEN 1990 AND 2023 AND "name" = 'Max' ORDER BY "year" DESC OFFSET 200 LIMIT 100

Generated SQL with placeholders request: 
SELECT "id", "name", "year" FROM "example_table" WHERE "year" BETWEEN $1 AND $2 AND "name" = $3 ORDER BY "year" DESC OFFSET 200 LIMIT 100

List of arguments for placeholders: 
"1990" "2023" "Max"

```
## Update [example.go](examples/pg/update.go)
### Code
```go
...
func update() {
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
```
### Output
```
Generated SQL plain request: 
UPDATE "example_table" SET "name" = 'John', "updated_at" = '2023-05-08 18:13:41.062736+03:00:00' WHERE "name" = 'Jjjohn' RETURNING "id"

Generated SQL with placeholders request: 
UPDATE "example_table" SET "name" = $1, "updated_at" = $2 WHERE "name" = $3 RETURNING "id"

List of arguments for placeholders: 
"John" "2023-05-08 18:13:41.062736 +0300 MSK m=+0.000647126" "Jjjohn"

```
## Delete [example.go](examples/pg/delete.go)
### Code
```go
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
```
### Output
```
Generated SQL plain request: 
DELETE FROM "example_table" WHERE "year" != 1990 RETURNING "id"

Generated SQL with placeholders request: 
DELETE FROM "example_table" WHERE "year" != $1 RETURNING "id"

List of arguments for placeholders: 
"1990"
```
## Insert [example.go](examples/pg/insert.go)
### Code
```go
func insert() {
    qb := pg.NewQueryBuilder()
    prepQuery := qb.InsertInto("example_table").
        Set("name", "John").
        Set("year", 1989)
    
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
```
### Output
```
Generated SQL plain request: 
INSERT INTO "example_table" ("name", "year") VALUES ('John', 1989)

Generated SQL with placeholders request: 
INSERT INTO "example_table" ("name", "year") VALUES ($1, $2)

List of arguments for placeholders: 
"John" "1989"

```