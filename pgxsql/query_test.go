package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
)

type rowsT struct {
}

func (r *rowsT) Close()     {}
func (r *rowsT) Err() error { return nil }
func (r *rowsT) CommandTag() CommandTag {
	return CommandTag{Sql: "select *", RowsAffected: 1, Insert: false, Update: false, Delete: false, Select: true}
}
func (r *rowsT) FieldDescriptions() []FieldDescription { return nil }
func (r *rowsT) Next() bool                            { return false }
func (r *rowsT) Scan(dest ...any) error                { return nil }
func (r *rowsT) Values() ([]any, error)                { return nil, nil }
func (r *rowsT) RawValues() [][]byte                   { return nil }

const (
	queryErrorSql  = "select * from test"
	queryRowsSql   = "select * from table"
	queryErrorPath = "query.error"
	queryRowsPath  = "query.rows"
)

func queryTestProxy(req Request) (Rows, error) {
	switch req.Uri {
	case BuildUri(queryErrorPath):
		return nil, errors.New("sqldml query error")
	case BuildUri(queryRowsPath):
		var i Rows = &rowsT{}
		return i, nil
	}
	return nil, nil
}

var qCtx = ContextWithQuery(context.Background(), queryTestProxy)

func ExampleQuery_Error() {
	rows, status := Query[template.DebugError](qCtx, NewRequest(queryErrorPath, queryErrorSql))
	fmt.Printf("test: Query[template.DebugError](ctx,%v) -> [rows:%v] [status:%v] \n", queryErrorSql, rows, status)

	//Output:
	//[[] github.com/idiomatic-go/postgresql-adapter/pgxsql/exec [sqldml query error]]
	//test: Query[template.DebugError](ctx,select * from test) -> [rows:<nil>] [status:Internal]
	
}

func ExampleQuery_Rows() {
	rows, status := Query[template.DebugError](qCtx, NewRequest(queryRowsPath, queryRowsSql))
	fmt.Printf("test: Query[template.DebugError](ctx,%v) -> [rows:%v] [status:%v] [cmd:%v]\n", queryRowsSql, rows, status, rows.CommandTag())

	//Output:
	//test: Query[template.DebugError](ctx,select * from table) -> [rows:&{}] [status:OK] [cmd:{select * 1 false false false true}]

}
