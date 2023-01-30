package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
)

func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

const (
	execUpdateSql  = "update test"
	execInsertSql  = "insert test"
	execUpdatePath = "exec.update"
	execInsertPath = "exec.insert"
)

func execTestProxy(req Request) (CommandTag, error) {
	switch req.Uri {
	case BuildUri(execUpdatePath):
		return emptyCommandTag, errors.New("exec error")
	case BuildUri(execInsertPath):
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

func ExampleExec() {
	ctx := ContextWithExec(context.Background(), execTestProxy)

	cmd, status := Exec[template.DebugError](ctx, NewRequest(execUpdatePath, execUpdateSql))
	fmt.Printf("test: Exec(%v) -> %v [cmd:%v]\n", execUpdateSql, status, cmd)

	cmd, status = Exec[template.DebugError](ctx, NewRequest(execInsertPath, execInsertSql))
	fmt.Printf("test: Exec(%v) -> %v [cmd:%v]\n", execInsertSql, status, cmd)

	//Output:
	//[[] github.com/idiomatic-go/postgresql-adapter/pgxsql/exec [exec error]]
	//test: Exec(update test) -> Internal [cmd:{ 0 false false false false}]
	//test: Exec(insert test) -> OK [cmd:{insert test 1234 true false false false}]

}
