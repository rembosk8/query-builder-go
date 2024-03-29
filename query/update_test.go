package query_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rembosk8/query-builder-go/query"
)

//nolint:funlen
func TestPGUpdate(t *testing.T) {
	const tableName = "tableName"

	t.Parallel()

	t.Run("update one ident for all records", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go'", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 1)
		assert.Equal(t, "go", args[0])
	})

	t.Run("update multiple fields for all records", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Set("year", 1989).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go', \"year\" = 1989", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Set("year", 1989).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1, \"year\" = $2", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
	})

	t.Run("update one ident WHERE", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Where("year").Equal(1989).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go' WHERE \"year\" = 1989", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Where("year").Equal(1989).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1 WHERE \"year\" = $2", tableName)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
		assert.NoError(t, err)
	})

	t.Run("update only", func(t *testing.T) {
		sql, err := qb.Update(tableName).Only().Set("name", "go").ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE ONLY %q SET \"name\" = 'go'", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Only().Set("name", "go").ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE ONLY %q SET \"name\" = $1", tableName)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "go")
		assert.NoError(t, err)
	})

	t.Run("update to DEFAULT", func(t *testing.T) {
		sql, err := qb.Update(tableName).Only().Set("name", "DEFAULT").Set("year", 1990).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE ONLY %q SET \"name\" = DEFAULT, \"year\" = 1990", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Only().Set("name", "DEFAULT").Set("year", 1990).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE ONLY %q SET \"name\" = DEFAULT, \"year\" = $1", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], 1990)
	})

	t.Run("update one ident for all records with RETURNING", func(t *testing.T) {
		sql, err := qb.Update(tableName).Set("name", "go").Returning("id", "name").ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go' RETURNING \"id\", \"name\"", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Set("name", "go").Returning("id", "name").ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1 RETURNING \"id\", \"name\"", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 1)
		assert.Equal(t, args[0], "go")
	})

	t.Run("update with multiple sets", func(t *testing.T) {
		sql, err := qb.Update(tableName).Sets("name", "go", "year", 1990).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go', \"year\" = 1990", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Sets("name", "go", "year", 1990).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1, \"year\" = $2", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
	})

	t.Run("update with multiple sets with FldVal", func(t *testing.T) {
		sql, err := qb.Update(tableName).Sets(query.FldVal("name", "go"), query.FldVal("year", 1990)).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go', \"year\" = 1990", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Sets(query.FldVal("name", "go"), query.FldVal("year", 1990)).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1, \"year\" = $2", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
	})

	t.Run("update with multiple sets with combined approach", func(t *testing.T) {
		sql, err := qb.Update(tableName).Sets(query.FldVal("name", "go"), "year", 1990).ToSQL()
		expectedSQL := fmt.Sprintf("UPDATE %q SET \"name\" = 'go', \"year\" = 1990", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Update(tableName).Sets(query.FldVal("name", "go"), "year", 1990).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("UPDATE %q SET \"name\" = $1, \"year\" = $2", tableName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSQL, sql)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
	})
}
