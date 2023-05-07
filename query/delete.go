package query

type Delete struct {
	baseQuery
}

var _ sqler = &Delete{}

func (d Delete) ToSql() (sql string, err error) {
	//TODO implement me
	panic("implement me")
}

func (d Delete) ToSqlWithStmts() (sql string, args []any, err error) {
	//TODO implement me
	panic("implement me")
}

func (d Delete) buildSqlPlain() {
	//TODO implement me
	panic("implement me")
}

func (d Delete) buildPrepStatement() (args []any) {
	//TODO implement me
	panic("implement me")
}
