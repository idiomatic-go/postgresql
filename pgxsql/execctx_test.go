package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/motif/template"
)

const (
	execTestUpdateSql = "update test"
	execTestInsertSql = "insert test"
	execTestUpdateRsc = "update"
	execTestInsertRsc = "insert"
)

var execCtxExchange = NewExecExchange(execCtxProxy)

func execCtxProxy(req *Request) (tag CommandTag, err error) {
	switch req.Uri {
	case BuildUpdateUri(execTestUpdateRsc):
		return tag, errors.New("exec error")
	case BuildInsertUri(execTestInsertRsc):
		return CommandTag{
			Sql:          "INSERT 1",
			RowsAffected: 1234,
			Insert:       true,
			Update:       false,
			Delete:       false,
			Select:       false,
		}, nil
	}
	return tag, nil
}

func ExampleExecContext_Error() {
	req := NewUpdateRequest(execTestUpdateRsc, execTestUpdateSql, nil, nil)
	tag, err := Exec[template.DebugError](context.Background(), NullCount, req)
	fmt.Printf("test: Exec[DebugError](ctx,NullCount,update) : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//[[] github.com/idiomatic-go/postgresql/pgxsql/exec [error on PostgreSQL exec call : dbClient is nil]]
	//test: Exec[DebugError](ctx,NullCount,update) : [tag:{ 0 false false false false}] [error:InvalidArgument]

}

func ExampleExecContext_Insert() {
	ctx := NewExecContext(nil, execCtxExchange)
	req := NewInsertRequest(execTestInsertRsc, execTestInsertSql, nil)
	tag, err := Exec[template.DebugError](ctx, NullCount, req)
	fmt.Printf("test: Exec[DebugError](ctx,NullCount,insert) : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//test: Exec[DebugError](ctx,NullCount,insert) : [tag:{INSERT 1 1234 true false false false}] [error:OK]

}
