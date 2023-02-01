package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type TestConditions struct {
	Time        time.Time
	Location    string
	Temperature float64
}

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
	queryErrorSql   = "select * from test"
	queryRowsSql    = "select * from table"
	queryConditions = "select * from conditions"
	queryErrorRsc   = "error"
	queryRowsRsc    = "rows"
)

func queryTestProxy(req Request) (Rows, error) {
	switch req.Uri {
	case BuildQueryUri(queryErrorRsc):
		return nil, errors.New("sqldml query error")
	case BuildQueryUri(queryRowsRsc):
		var i Rows = &rowsT{}
		return i, nil
	}
	return nil, nil
}

var qCtx = ContextWithQuery(context.Background(), queryTestProxy)

func ExampleQuery_TestError() {
	result, status := Query[template.DebugError](qCtx, NewQueryRequest(queryErrorRsc, queryErrorSql))
	fmt.Printf("test: Query[template.DebugError](ctx,%v) -> [rows:%v] [status:%v]\n", queryErrorSql, result, status)

	//Output:
	//[[] github.com/idiomatic-go/postgresql/pgxsql/exec [sqldml query error]]
	//test: Query[template.DebugError](ctx,select * from test) -> [rows:<nil>] [status:Internal]

}

func ExampleQuery_TestRows() {
	result, status := Query[template.DebugError](qCtx, NewQueryRequest(queryRowsRsc, queryRowsSql))
	fmt.Printf("test: Query[template.DebugError](ctx,%v) -> [rows:%v] [status:%v] [cmd:%v]\n", queryRowsSql, result, status, result.CommandTag())

	//Output:
	//test: Query[template.DebugError](ctx,select * from table) -> [rows:&{}] [status:OK] [cmd:{select * 1 false false false true}]

}

func ExampleQuery_Conditions() {
	req := NewQueryRequest(queryRowsRsc, queryConditions)

	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		results, status := Query[template.DebugError](nil, req)
		if !status.OK() {
			fmt.Printf("test: Query[template.DebugError](nil,%v) -> [status:%v]\n", queryConditions, status)
		} else {
			fmt.Printf("test: Query[template.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", queryConditions, status, results.CommandTag())
			conditions, status1 := processResults(results, "")
			fmt.Printf("test: processResults(results) -> [status:%v] [rows:%v]\n", status1, conditions)
		}
	}

	//Output:
	//test: Query[template.DebugError](nil,select * from conditions) -> [status:OK] [cmd:{ 0 false false false false}]
	//test: processResults(results) -> [status:OK] [rows:[{2023-01-26 12:09:12.426535 -0600 CST office 70} {2023-01-26 12:09:12.426535 -0600 CST basement 66.5} {2023-01-26 12:09:12.426535 -0600 CST garage 77}]]

}

func processResults(results Rows, msg string) (conditions []TestConditions, status *template.Status) {
	conditions, status = scanRows(results)
	if status.OK() && len(conditions) == 0 {
		return nil, template.NewStatusCode(template.StatusNotFound)
	}
	return conditions, status
}

func scanRows(rows Rows) ([]TestConditions, *template.Status) {
	if rows == nil {
		return nil, template.NewStatusInvalidArgument("", errors.New("invalid request: Rows interface is empty"))
	}
	var err error
	var values []any
	var conditions []TestConditions
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, template.NewStatusError("", err)
		}
		values, err = rows.Values()
		if err != nil {
			return nil, template.NewStatusError("", err)
		}
		conditions = append(conditions, scanColumns(values))
	}
	return conditions, template.NewStatusOK()
}

func scanColumns(values []any) TestConditions {
	var ts = new(pgtype.Timestamp)
	ts.Scan(values[0])

	cond := TestConditions{
		Time:        ts.Time,
		Location:    values[1].(string),
		Temperature: values[2].(float64),
	}
	return cond
}
