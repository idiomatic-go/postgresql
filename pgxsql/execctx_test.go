package pgxsql

import (
	"context"
	"errors"
	"fmt"
)

const (
	execTestUpdateSql = "update test"
	execTestInsertSql = "insert test"
	execTestUpdateRsc = "update"
	execTestInsertRsc = "insert"
)

func execCtxProxy(req Request) (CommandTag, error) {
	switch req.Uri {
	case BuildExecUri(execTestUpdateRsc):
		return emptyCommandTag, errors.New("exec error")
	case BuildExecUri(execTestInsertRsc):
		return CommandTag{
			Sql:          req.Sql,
			RowsAffected: 1234,
			Insert:       true,
			Update:       false,
			Delete:       false,
			Select:       false,
		}, nil
	}
	return emptyCommandTag, nil
}

func ExampleContextExec_Error() {
	ctx := ContextWithExec(context.Background(), execCtxProxy)
	tag, err := ContextExec(ctx, NewExecRequest(execTestUpdateRsc, execTestUpdateSql))
	fmt.Printf("test: ExecQuery() : [tags:%v] [error:%v]\n", tag, err)

	//Output:
	//test: ExecQuery() : [tags:{ 0 false false false false}] [error:exec error]

}

func ExampleContextExec_Rows() {
	ctx := ContextWithExec(context.Background(), execCtxProxy)
	tag, err := ContextExec(ctx, NewExecRequest(execTestInsertRsc, execTestInsertSql))
	fmt.Printf("test: ContextExec() : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//test: ContextExec() : [tag:{insert test 1234 true false false false}] [error:<nil>]

}
