package pgxsql

import (
	"context"
	"errors"
	"fmt"
)

const (
	queryTestErrorSql = "select * from test"
	queryTestRowsSql  = "select * from table"
	queryTestErrorRsc = "error"
	queryTestRowsRsc  = "rows"
)

func queryctxProxy(req *Request) (Rows, error) {
	switch req.Uri {
	case BuildQueryUri(queryTestErrorRsc):
		return nil, errors.New("pgxsql query error")
	case BuildQueryUri(queryTestRowsRsc):
		var i Rows = &rowsT{}
		return i, nil
	}
	return nil, nil
}

func ExampleContextQuery_Error() {
	ctx := ContextWithQuery(context.Background(), queryctxProxy)
	rows, err := ContextQuery(ctx, NewQueryRequest(queryTestErrorRsc, queryTestErrorSql, nil))
	fmt.Printf("test: ContextQuery() : [rows:%v] [error:%v]\n", rows != nil, err)

	//Output:
	//test: ContextQuery() : [rows:false] [error:pgxsql query error]

}

func ExampleContextQuery_Rows() {
	ctx := ContextWithQuery(context.Background(), queryctxProxy)
	rows, err := ContextQuery(ctx, NewQueryRequest(queryTestRowsRsc, queryTestRowsSql, nil))
	fmt.Printf("test: ContextQuery() : [rows:%v] [error:%v]\n", rows != nil, err)

	//Output:
	//test: ContextQuery() : [rows:true] [error:<nil>]

}
