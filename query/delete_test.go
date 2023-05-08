package query_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGDelete(t *testing.T) {
	tableName := "tableName"

	t.Parallel()

	t.Run("delete all records from table", func(t *testing.T) {
		sql, args, err := qb.DeleteFrom(tableName).ToSqlWithStmts()
		expectedSql := fmt.Sprintf("DELETE FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 0)
		assert.NoError(t, err)

		sql, err = qb.DeleteFrom(tableName).ToSql()
		expectedSql = fmt.Sprintf("DELETE FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)
	})

	t.Run("delete all records from ONLY table", func(t *testing.T) {
		sql, args, err := qb.DeleteFrom(tableName).Only().ToSqlWithStmts()
		expectedSql := fmt.Sprintf("DELETE FROM ONLY \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 0)
		assert.NoError(t, err)

		sql, err = qb.DeleteFrom(tableName).Only().ToSql()
		expectedSql = fmt.Sprintf("DELETE FROM ONLY \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)
	})

	t.Run("delete some recs from table WHERE", func(t *testing.T) {
		sql, args, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").ToSqlWithStmts()
		expectedSql := fmt.Sprintf("DELETE FROM \"%s\" WHERE \"name\" <= $1", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "AAA")

		sql, err = qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").ToSql()
		expectedSql = fmt.Sprintf("DELETE FROM \"%s\" WHERE \"name\" <= 'AAA'", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)
	})

	t.Run("delete some recs from table WHERE and RETURNING", func(t *testing.T) {
		sql, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").Returning("id", "name").ToSql()
		expectedSql := fmt.Sprintf("DELETE FROM \"%s\" WHERE \"name\" <= 'AAA' RETURNING \"id\", \"name\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").Returning("id", "name").ToSqlWithStmts()
		expectedSql = fmt.Sprintf("DELETE FROM \"%s\" WHERE \"name\" <= $1 RETURNING \"id\", \"name\"", tableName)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "AAA")
		assert.NoError(t, err)
	})
}
