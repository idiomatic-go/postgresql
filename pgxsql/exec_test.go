package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"github.com/idiomatic-go/postgresql/pgxdml"
	"time"
)

func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

const (
	execUpdateSql = "update test"
	execInsertSql = "insert test"
	execUpdateRsc = "update"
	execInsertRsc = "insert"
	execDeleteRsc = "insert"

	execInsertConditions = "INSERT INTO conditions (time,location,temperature) VALUES"
	execUpdateConditions = "UPDATE conditions"
	execDeleteConditions = "DELETE FROM conditions"
)

func execTestProxy(req Request) (CommandTag, error) {
	switch req.Uri {
	case BuildExecUri(execUpdateRsc):
		return emptyCommandTag, errors.New("exec error")
	case BuildExecUri(execInsertRsc):
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

	cmd, status := Exec[template.DebugError](ctx, NullCount, NewExecRequest(execUpdateRsc, execUpdateSql))
	fmt.Printf("test: Exec(%v) -> %v [cmd:%v]\n", execUpdateSql, status, cmd)

	cmd, status = Exec[template.DebugError](ctx, NullCount, NewExecRequest(execInsertRsc, execInsertSql))
	fmt.Printf("test: Exec(%v) -> %v [cmd:%v]\n", execInsertSql, status, cmd)

	//Output:
	//[[] github.com/idiomatic-go/postgresql/pgxsql/exec [exec error]]
	//test: Exec(update test) -> Internal [cmd:{ 0 false false false false}]
	//test: Exec(insert test) -> OK [cmd:{insert test 1234 true false false false}]

}

func ExampleExec_Insert() {
	req := NewExecRequest(execInsertRsc, execInsertConditions)
	cond := TestConditions{
		Time:        time.Now(),
		Location:    "frisco",
		Temperature: 27.66,
	}
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		//values := []any{pgxdml.Function(pgxdml.ChangedTimestampFn), cond.Location, cond.Temperature}
		//stmt, err := pgxdml.WriteInsert(req.Sql, values)
		//fmt.Printf("test: WriteInsert() -> [error:%v] [sql:%v}\n", err, stmt)

		results, status := Insert[template.DebugError](nil, NullCount, req, pgxdml.NewInsertValues([]any{pgxdml.TimestampFn, cond.Location, cond.Temperature}))
		if !status.OK() {
			fmt.Printf("test: ExecInsert[template.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execInsertConditions, status, results)
		} else {
			fmt.Printf("test: ExecInsert[template.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execInsertConditions, status, results)
		}
	}

	//Output:
	//test: ExecInsert[template.DebugError](nil,INSERT INTO conditions (time,location,temperature) VALUES) -> [status:OK] [cmd:{INSERT 0 1 1 true false false false}]

}

func ExampleExec_Update() {
	req := NewExecRequest(execUpdateRsc, execUpdateConditions)

	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		attrs := []pgxdml.Attr{{"Temperature", 45.1234}}
		where := []pgxdml.Attr{{"Location", "garage"}}

		//stmt, err := pgxdml.WriteUpdate(req.Sql, attrs, where)
		//fmt.Printf("test: WriteUpdate() -> [error:%v] [sql:%v]\n", err, stmt)

		results, status := Update[template.DebugError](nil, NullCount, req, attrs, where)
		if !status.OK() {
			fmt.Printf("test: ExecUpdate[template.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execUpdateConditions, status, results)
		} else {
			fmt.Printf("test: ExecUpdate[template.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execUpdateConditions, status, results)
		}
	}

	//Output:
	//test: ExecUpdate[template.DebugError](nil,UPDATE conditions) -> [status:OK] [cmd:{UPDATE 1 1 false true false false}]

}

func ExampleExec_Delete() {
	req := NewExecRequest(execDeleteRsc, execDeleteConditions)

	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		where := []pgxdml.Attr{{"Location", "frisco"}}
		//stmt, err := pgxdml.WriteDelete(req.Sql, where)
		//fmt.Printf("test: WriteDelete() -> [error:%v] [sql:%v]\n", err, stmt)

		results, status := Delete[template.DebugError](nil, NullCount, req, where)
		if !status.OK() {
			fmt.Printf("test: ExecDelete[template.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execDeleteConditions, status, results)
		} else {
			fmt.Printf("test: ExecDelete[template.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execDeleteConditions, status, results)
		}
	}

	//Output:
	//test: ExecDelete[template.DebugError](nil,DELETE FROM conditions) -> [status:OK] [cmd:{DELETE 1 1 false false true false}]

}
