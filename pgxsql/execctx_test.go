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

func execCtxProxy(req *Request) (CommandTag, error) {
	switch req.Uri {
	case BuildUpdateUri(execTestUpdateRsc):
		return emptyCommandTag, errors.New("exec error")
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
	return emptyCommandTag, nil
}

func ExampleContextExec_Error() {
	ctx := ContextWithExec(context.Background(), execCtxProxy)
	tag, err := ContextExec(ctx, NewUpdateRequest(execTestUpdateRsc, execTestUpdateSql, nil, nil))
	fmt.Printf("test: ExecQuery() : [tags:%v] [error:%v]\n", tag, err)

	//Output:
	//test: ExecQuery() : [tags:{ 0 false false false false}] [error:exec error]

}

func ExampleContextExec_Rows() {
	ctx := ContextWithExec(context.Background(), execCtxProxy)
	tag, err := ContextExec(ctx, NewInsertRequest(execTestInsertRsc, execTestInsertSql, nil))
	fmt.Printf("test: ContextExec() : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//test: ContextExec() : [tag:{INSERT 1 1234 true false false false}] [error:<nil>]

}
