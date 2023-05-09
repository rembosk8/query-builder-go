package query_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGUpdate(t *testing.T) {
	tableName := "tableName"

	t.Parallel()

	t.Run("update one indend for all records", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").ToSql()
		expectedSql := fmt.Sprintf("UPDATE \"%s\" SET \"name\" = 'go'", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE \"%s\" SET \"name\" = $1", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "go")
	})

	t.Run("update multiple fields for all records", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Set("year", 1989).ToSql()
		expectedSql := fmt.Sprintf("UPDATE \"%s\" SET \"name\" = 'go', \"year\" = 1989", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Set("year", 1989).ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE \"%s\" SET \"name\" = $1, \"year\" = $2", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
	})

	t.Run("update one indend WHERE", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Where("year").Equal(1989).ToSql()
		expectedSql := fmt.Sprintf("UPDATE \"%s\" SET \"name\" = 'go' WHERE \"year\" = 1989", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Where("year").Equal(1989).ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE \"%s\" SET \"name\" = $1 WHERE \"year\" = $2", tableName)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
		assert.NoError(t, err)
	})

	t.Run("update only", func(t *testing.T) {
		sql, err := qb.Update(tableName).Only().Set("name", "go").ToSql()
		expectedSql := fmt.Sprintf("UPDATE ONLY \"%s\" SET \"name\" = 'go'", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Only().Set("name", "go").ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE ONLY \"%s\" SET \"name\" = $1", tableName)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "go")
		assert.NoError(t, err)
	})

	t.Run("update to DEFAULT", func(t *testing.T) {
		sql, err := qb.Update(tableName).Only().Set("name", "DEFAULT").Set("year", 1990).ToSql()
		expectedSql := fmt.Sprintf("UPDATE ONLY \"%s\" SET \"name\" = DEFAULT, \"year\" = 1990", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Only().Set("name", "DEFAULT").Set("year", 1990).ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE ONLY \"%s\" SET \"name\" = DEFAULT, \"year\" = $1", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], 1990)
	})

	t.Run("update one indend for all records with RETURNING", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Returning("id", "name").ToSql()
		expectedSql := fmt.Sprintf("UPDATE \"%s\" SET \"name\" = 'go' RETURNING \"id\", \"name\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Returning("id", "name").ToSqlWithStmts()
		expectedSql = fmt.Sprintf("UPDATE \"%s\" SET \"name\" = $1 RETURNING \"id\", \"name\"", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "go")
	})
}
