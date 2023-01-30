package pgxsql

import (
	"context"
	"errors"
	"fmt"
)

const (
	queryTestErrorSql  = "select * from test"
	queryTestRowsSql   = "select * from table"
	queryTestErrorPath = "query.error"
	queryTestRowsPath  = "query.rows"
)

func queryctxProxy(req Request) (Rows, error) {
	switch req.Uri {
	case BuildUri(queryTestErrorPath):
		return nil, errors.New("sqldml query error")
	case BuildUri(queryTestRowsPath):
		var i Rows = &rowsT{}
		return i, nil
	}
	return nil, nil
}

func ExampleContextQuery_Error() {
	ctx := ContextWithQuery(context.Background(), queryctxProxy)
	rows, err := ContextQuery(ctx, NewRequest(queryTestErrorPath, queryTestErrorSql))
	fmt.Printf("test: ContextQuery() : [rows:%v] [error:%v]\n", rows != nil, err)

	//Output:
	//test: ContextQuery() : [rows:false] [error:sqldml query error]

}

func ExampleContextQuery_Rows() {
	ctx := ContextWithQuery(context.Background(), queryctxProxy)
	rows, err := ContextQuery(ctx, NewRequest(queryTestRowsPath, queryTestRowsSql))
	fmt.Printf("test: ContextQuery() : [rows:%v] [error:%v]\n", rows != nil, err)

	//Output:
	//test: ContextQuery() : [rows:true] [error:<nil>]

}
