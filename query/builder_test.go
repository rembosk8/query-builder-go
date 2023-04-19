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
}
