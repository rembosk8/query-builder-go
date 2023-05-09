package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhereAll(t *testing.T) {
	prepQuery := qb.Select().From("example_table")

	prepQueryWithWhere := prepQuery.Where("name").Equal("Max").
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

	expectedSql := "SELECT * FROM \"example_table\" WHERE \"name\" = 'Max' AND \"last_name\" != 'Brown' AND \"year\" BETWEEN 1990 AND 2023 AND \"year\" NOT BETWEEN 2000 AND 2010 AND \"salary\" < 2000 AND \"salary\" > 1000 AND \"days\" >= 10 AND \"days\" <= 30 AND \"department\" IN ('HR', 'Development') AND \"position\" NOT IN ('junior', 'middle') AND \"referral\" IS NULL AND \"lead\" IS NOT NULL AND \"comment\" LIKE '%something interesting%' AND \"comment\" NOT LIKE '%something not interesting%'"
	sql, err := prepQueryWithWhere.ToSql()
	require.NoError(t, err)
	assert.Equal(t, expectedSql, sql)

	expectedSql = "SELECT * FROM \"example_table\" WHERE \"name\" = $1 AND \"last_name\" != $2 AND \"year\" BETWEEN $3 AND $4 AND \"year\" NOT BETWEEN $5 AND $6 AND \"salary\" < $7 AND \"salary\" > $8 AND \"days\" >= $9 AND \"days\" <= $10 AND \"department\" IN ($11, $12) AND \"position\" NOT IN ($13, $14) AND \"referral\" IS NULL AND \"lead\" IS NOT NULL AND \"comment\" LIKE $15 AND \"comment\" NOT LIKE $16"
	sql, args, err := prepQueryWithWhere.ToSqlWithStmts()
	require.NoError(t, err)
	assert.Equal(t, expectedSql, sql)
	assert.Len(t, args, 16)
}
