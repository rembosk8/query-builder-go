package query_test

import (
	"fmt"
	"testing"

	"github.com/rembosk8/query-builder-go/query/pg"
	"github.com/stretchr/testify/assert"
)

func TestPGQueryBuilder(t *testing.T) {
	qb := pg.NewQueryBuilder()
	tableName := "tableName"

	t.Parallel()

	t.Run("not initialized builder", func(t *testing.T) {
		sql, args, err := qb.Build()
		assert.Empty(t, sql, "must return empty SQL when nothing is initialized")
		assert.Nil(t, args, "must return nil for args when nothing is initialized")
		assert.Error(t, err, "must return err")
	})

	t.Run("try to select all from specified table", func(t *testing.T) {
		sql, args, err := qb.From(tableName).Build()
		expectedSql := fmt.Sprintf("SELECT * FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select id from specified table", func(t *testing.T) {
		sql, args, err := qb.Select("id").From(tableName).Build()
		expectedSql := fmt.Sprintf("SELECT \"id\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select multiple fields from specified table", func(t *testing.T) {
		sql, args, err := qb.Select("id", "name", "value").From(tableName).Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select multiple fields from specified table with WHERE", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").Equal(1).Where("name").Equal("testName")
		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName'", tableName)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = $1 AND \"name\" = $2", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 2)
		assert.Equal(t, 1, args[0])
		assert.Equal(t, "testName", args[1])
		assert.NoError(t, err)
	})

	t.Run("select with TOP and LIMIT", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").Equal(1).
			Where("name").Equal("testName").
			Offset(10).
			Limit(5)
		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName' OFFSET 10 LIMIT 5", tableName)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = $1 AND \"name\" = $2 OFFSET 10 LIMIT 5", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 2)
		assert.Equal(t, 1, args[0])
		assert.Equal(t, "testName", args[1])
		assert.NoError(t, err)
	})

	t.Run("select with TOP and LIMIT and ORDER BY", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").Equal(1).
			Where("name").Equal("testName").
			Offset(10).
			Limit(5).
			OrderBy("id").Desc().
			OrderBy("name").Asc()
		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName' ORDER BY \"id\" DESC, \"name\" ASC OFFSET 10 LIMIT 5",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = $1 AND \"name\" = $2 ORDER BY \"id\" DESC, \"name\" ASC OFFSET 10 LIMIT 5", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 2)
		assert.Equal(t, 1, args[0])
		assert.Equal(t, "testName", args[1])
		assert.NoError(t, err)
	})

	t.Run("select where IN condition", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").In(1, 2, 3).
			Where("name").NotIn("name 1", "name 2")

		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IN (1, 2, 3) AND \"name\" NOT IN ('name 1', 'name 2')",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IN ($1, $2, $3) AND \"name\" NOT IN ($4, $5)", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 5)
		assert.Equal(t, 1, args[0])
		assert.Equal(t, 2, args[1])
		assert.Equal(t, 3, args[2])
		assert.Equal(t, "name 1", args[3])
		assert.Equal(t, "name 2", args[4])
		assert.NoError(t, err)
	})

	t.Run("select where IS NULL and IS NOT NULL condition", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").IsNull().
			Where("name").IsNotNull()

		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IS NULL AND \"name\" IS NOT NULL",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IS NULL AND \"name\" IS NOT NULL", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 0)
		assert.NoError(t, err)
	})

	t.Run("select where BETWEEN and NOT BETWEEN", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").Between(1, 10).
			Where("name").NotBetween("first", "second")

		sql, err := prebuild.BuildPlain()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" BETWEEN 1 AND 10 AND \"name\" NOT BETWEEN 'first' AND 'second'",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.Build()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" BETWEEN $1 AND $2 AND \"name\" NOT BETWEEN $3 AND $4", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 4)
		assert.Equal(t, 1, args[0])
		assert.Equal(t, 10, args[1])
		assert.Equal(t, "first", args[2])
		assert.Equal(t, "second", args[3])
		assert.NoError(t, err)
	})
}
