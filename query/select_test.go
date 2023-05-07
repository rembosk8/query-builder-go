package query_test

import (
	"fmt"
	"testing"

	"github.com/rembosk8/query-builder-go/query"
	"github.com/rembosk8/query-builder-go/query/pg"
	"github.com/stretchr/testify/assert"
)

func TestPGSelect(t *testing.T) {
	tableName := "tableName"

	t.Parallel()

	t.Run("try to select all from specified table", func(t *testing.T) {
		sql, err := qb.Select().From(tableName).ToSql()
		expectedSql := fmt.Sprintf("SELECT * FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.NoError(t, err)

		sql, args, err := qb.Select().From(tableName).ToSqlWithStmts()
		expectedSql = fmt.Sprintf("SELECT * FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select id from specified table", func(t *testing.T) {
		sql, args, err := qb.Select("id").From(tableName).ToSqlWithStmts()
		expectedSql := fmt.Sprintf("SELECT \"id\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)

		sql, err = qb.Select("id").From(tableName).ToSql()
		expectedSql = fmt.Sprintf("SELECT \"id\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select multiple fields from specified table", func(t *testing.T) {
		sql, args, err := qb.Select("id", "name", "value").From(tableName).ToSqlWithStmts()
		expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)

		sql, err = qb.Select("id", "name", "value").From(tableName).ToSql()
		expectedSql = fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\"", tableName)
		assert.Equal(t, expectedSql, sql)
		assert.Empty(t, args, "args must be empty")
		assert.NoError(t, err)
	})

	t.Run("try to select multiple fields from specified table with WHERE", func(t *testing.T) {
		prebuild := qb.Select("id", "name", "value").
			From(tableName).
			Where("id").Equal(1).Where("name").Equal("testName")
		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName'", tableName)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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
		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf("SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName' OFFSET 10 LIMIT 5", tableName)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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
		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" = 1 AND \"name\" = 'testName' ORDER BY \"id\" DESC, \"name\" ASC OFFSET 10 LIMIT 5",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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

		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IN (1, 2, 3) AND \"name\" NOT IN ('name 1', 'name 2')",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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

		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" IS NULL AND \"name\" IS NOT NULL",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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

		sql, err := prebuild.ToSql()
		expectedPlainSql := fmt.Sprintf(
			"SELECT \"id\", \"name\", \"value\" FROM \"%s\" WHERE \"id\" BETWEEN 1 AND 10 AND \"name\" NOT BETWEEN 'first' AND 'second'",
			tableName,
		)
		assert.Equal(t, expectedPlainSql, sql)
		assert.NoError(t, err)

		sql, args, err := prebuild.ToSqlWithStmts()
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

func TestSelectReusage(t *testing.T) {
	var (
		sql  string
		err  error
		args []any
	)
	tableName := "tableName"

	prepBuild := qb.Select("id", "name", "year").From(tableName)

	sql, args, err = prepBuild.Where("first").Equal(1).ToSqlWithStmts()
	expectedSql := fmt.Sprintf("SELECT \"id\", \"name\", \"year\" FROM \"%s\" WHERE \"first\" = $1", tableName)
	assert.Equal(t, expectedSql, sql)
	assert.NoError(t, err)
	assert.Len(t, args, 1)
	assert.Equal(t, 1, args[0])

	sql, args, err = prepBuild.Where("first2").Equal(10).Where("second").Equal(20).ToSqlWithStmts()
	expectedSql = fmt.Sprintf("SELECT \"id\", \"name\", \"year\" FROM \"%s\" WHERE \"first2\" = $1 AND \"second\" = $2", tableName)
	assert.Equal(t, expectedSql, sql)
	assert.NoError(t, err)
	assert.Len(t, args, 2)
	assert.Equal(t, 10, args[0])
	assert.Equal(t, 20, args[1])
}

func TestSelectV2CustomTag(t *testing.T) {
	tableName := "tableName"
	qb := query.New(query.WithIndentBuilder(pg.IndentBuilder()), query.WithStructAnnotationTag("myTag"))

	t.Run("select v2 with custom field tag", func(t *testing.T) {
		type tableModelWithAnnotation struct {
			ID   string `myTag:"id_a"`
			Name string `myTag:"name_a"`
		}
		m := tableModelWithAnnotation{}
		sql, args, err := qb.SelectV2(&m).From(tableName).ToSqlWithStmts()
		expectedSql := `SELECT "id_a", "name_a" FROM "tableName"`
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 0)
	})
	t.Run("select v2 without tags", func(t *testing.T) {
		type tableModel struct {
			ID       string
			LastName string
		}

		m2 := tableModel{}
		sql, args, err := qb.SelectV2(&m2).From(tableName).ToSqlWithStmts()
		expectedSql := `SELECT "id", "last_name" FROM "tableName"`
		assert.NoError(t, err)
		assert.Equal(t, expectedSql, sql)
		assert.Len(t, args, 0)
	})
}
