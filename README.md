# query-builder-go
SQL Query Builder for golang which helps to build SQL queries dynamically.
The goal for the first release is to support all common operations with rows for PostgreSQL, 
such as SELECT, INSERT, UPDATE, DELETE. 

Support of filtering with WHERE(AND), OFFSET, LIMIT, ORDER BY.

SQL Query Builder could provides 2 different result:
- Plain SQL query
```
INSERT INTO "example_table" ("name", "year") VALUES ('John', 1989)
```
- SQL query with placeholders + list of arguments for the placeholders.
```
Generated SQL with placeholders query: 
INSERT INTO "example_table" ("name", "year") VALUES ($1, $2)

List of arguments for placeholders: 
"John" "1989"
```

## Where [example](query/where_test.go)
Support the following operators:

```go
	prepQuery := qb.Select().From("example_table")

	prepQueryWithWhere := prepQuery.
		Where("name").Equal("Max").
		Where("last_name").NotEqual("Brown").
		Where("year").Between(1990, 2023).
		Where("year").NotBetween(2000, 2010).
		Where("salary").Less(2000).
		Where("salary").Greater(1000).
		Where("days").GreaterEqual(10).
		Where("days").LessEqual(30).
		Where("department").In("HR", "Development").
		Where("position").NotIn("junior", "middle").
		Where("referral").IsNull().
		Where("lead").IsNotNull().
		Where("comment").Like("%something interesting%").
		Where("comment").NotLike("%something not interesting%")
```

## Insert
Generate query for INSERT without any values is going to generate `INSERT INTO <tablename> DEFAULT VALUES` 

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

	fmt.Println("Generated SQL plain query: ")
	fmt.Println(sql)

	sql, args, err := prepQuery.ToSqlWithStmts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL with placeholders query: ")
	fmt.Println(sql)

	fmt.Println("List of arguments for placeholders: ")
	fmt.Println(args...)
}
```
### Output
```
Generated SQL plain query: 
SELECT "id", "name", "year" FROM "example_table" WHERE "year" BETWEEN 1990 AND 2023 AND "name" = 'Max' ORDER BY "year" DESC OFFSET 200 LIMIT 100

Generated SQL with placeholders query: 
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

	fmt.Println("Generated SQL plain query: ")
	fmt.Println(sql)

	sql, args, err := prepQuery.ToSqlWithStmts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated SQL with placeholders query: ")
	fmt.Println(sql)

	fmt.Println("List of arguments for placeholders: ")
	fmt.Println(args...)
}
```
### Output
```
Generated SQL plain query: 
UPDATE "example_table" SET "name" = 'John', "updated_at" = '2023-05-08 18:13:41.062736+03:00:00' WHERE "name" = 'Jjjohn' RETURNING "id"

Generated SQL with placeholders query: 
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
    
    fmt.Println("Generated SQL plain query: ")
    fmt.Println(sql)
    
    sql, args, err := prepQuery.ToSqlWithStmts()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Generated SQL with placeholders query: ")
    fmt.Println(sql)
    
    fmt.Println("List of arguments for placeholders: ")
    fmt.Println(args...)
}
```
### Output
```
Generated SQL plain query: 
DELETE FROM "example_table" WHERE "year" != 1990 RETURNING "id"

Generated SQL with placeholders query: 
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
    
    fmt.Println("Generated SQL plain query: ")
    fmt.Println(sql)
    
    sql, args, err := prepQuery.ToSqlWithStmts()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Generated SQL with placeholders query: ")
    fmt.Println(sql)
    
    fmt.Println("List of arguments for placeholders: ")
    fmt.Println(args...)
}
```
### Output
```
Generated SQL plain query: 
INSERT INTO "example_table" ("name", "year") VALUES ('John', 1989)

Generated SQL with placeholders query: 
INSERT INTO "example_table" ("name", "year") VALUES ($1, $2)

List of arguments for placeholders: 
"John" "1989"

```