package query_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGInsert(t *testing.T) {
	const tableName = "tableName"

	t.Parallel()

	t.Run("insert a record into the table", func(t *testing.T) {
		sql, err := qb.InsertInto(tableName).Set("name", "go").Set("year", 1989).ToSQL()
		expectedSQL := fmt.Sprintf("INSERT INTO %q (\"name\", \"year\") VALUES ('go', 1989)", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.InsertInto(tableName).Set("name", "go").Set("year", 1989).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("INSERT INTO %q (\"name\", \"year\") VALUES ($1, $2)", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
	})

	t.Run("insert a record into the table with sets", func(t *testing.T) {
		sql, err := qb.InsertInto(tableName).Sets("name", "go", "year", 1989).ToSQL()
		expectedSQL := fmt.Sprintf("INSERT INTO %q (\"name\", \"year\") VALUES ('go', 1989)", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.InsertInto(tableName).Sets("name", "go", "year", 1989).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("INSERT INTO %q (\"name\", \"year\") VALUES ($1, $2)", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)
		require.Len(t, args, 2)
		assert.Equal(t, args[0], "go")
		assert.Equal(t, args[1], 1989)
	})

	t.Run("insert DEFAULT VALUES", func(t *testing.T) {
		sql, err := qb.InsertInto(tableName).ToSQL()
		expectedSQL := fmt.Sprintf("INSERT INTO %q DEFAULT VALUES", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)

		sql, args, err := qb.InsertInto(tableName).ToSQLWithStmts()
		expectedSQL = fmt.Sprintf("INSERT INTO %q DEFAULT VALUES", tableName)
		assert.Equal(t, expectedSQL, sql)
		assert.NoError(t, err)
		require.Len(t, args, 0)
	})
}
