package query_test

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestPGDelete(t *testing.T) {
//	const tableName = "tableName"
//
//	t.Parallel()
//
//	t.Run("delete all records from table", func(t *testing.T) {
//		sql, args, err := qb.DeleteFrom(tableName).ToSQLWithStmts()
//		expectedSQL := fmt.Sprintf("DELETE FROM %q", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		require.Len(t, args, 0)
//		assert.NoError(t, err)
//
//		sql, err = qb.DeleteFrom(tableName).ToSQL()
//		expectedSQL = fmt.Sprintf("DELETE FROM %q", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		assert.NoError(t, err)
//	})
//
//	t.Run("delete all records from ONLY table", func(t *testing.T) {
//		sql, args, err := qb.DeleteFrom(tableName).Only().ToSQLWithStmts()
//		expectedSQL := fmt.Sprintf("DELETE FROM ONLY %q", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		require.Len(t, args, 0)
//		assert.NoError(t, err)
//
//		sql, err = qb.DeleteFrom(tableName).Only().ToSQL()
//		expectedSQL = fmt.Sprintf("DELETE FROM ONLY %q", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		assert.NoError(t, err)
//	})
//
//	t.Run("delete some recs from table WHERE", func(t *testing.T) {
//		sql, args, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").ToSQLWithStmts()
//		expectedSQL := fmt.Sprintf("DELETE FROM %q WHERE \"name\" <= $1", tableName)
//		assert.NoError(t, err)
//		assert.Equal(t, expectedSQL, sql)
//		require.Len(t, args, 1)
//		assert.Equal(t, args[0], "AAA")
//
//		sql, err = qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").ToSQL()
//		expectedSQL = fmt.Sprintf("DELETE FROM %q WHERE \"name\" <= 'AAA'", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		assert.NoError(t, err)
//	})
//
//	t.Run("delete some recs from table WHERE and RETURNING", func(t *testing.T) {
//		sql, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").Returning("id", "name").ToSQL()
//		expectedSQL := fmt.Sprintf("DELETE FROM %q WHERE \"name\" <= 'AAA' RETURNING \"id\", \"name\"", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		assert.NoError(t, err)
//
//		sql, args, err := qb.DeleteFrom(tableName).Where("name").LessEqual("AAA").Returning("id", "name").ToSQLWithStmts()
//		expectedSQL = fmt.Sprintf("DELETE FROM %q WHERE \"name\" <= $1 RETURNING \"id\", \"name\"", tableName)
//		assert.Equal(t, expectedSQL, sql)
//		require.Len(t, args, 1)
//		assert.Equal(t, args[0], "AAA")
//		assert.NoError(t, err)
//	})
//}
