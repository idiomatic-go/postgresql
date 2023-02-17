package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/runtime"
	"github.com/idiomatic-go/motif/template"
)

const (
	NullCount = int64(-1)
)

var execLoc = pkgPath + "/exec"

// Exec - templated function for executing a SQL statement
func Exec[E template.ErrorHandler](ctx context.Context, expectedCount int64, req *Request, args ...any) (_ CommandTag, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	fn, ctx, limited = actuatorApply(ctx, messaging.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return CommandTag{}, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if IsContextExec(ctx) {
		result, err := ContextExec(ctx, req)
		return result, e.HandleWithContext(ctx, execLoc, err)
	}
	if dbClient == nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err0)
	}
	t, err := dbClient.Exec(ctx, req.BuildSql(), args...)
	if err != nil {
		err0 = txn.Rollback(ctx)
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, recast(err), err0)
	}
	if expectedCount != NullCount && t.RowsAffected() != expectedCount {
		err0 = txn.Rollback(ctx)
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, errors.New(fmt.Sprintf("error exec statement [%v] : actual RowsAffected %v != expected RowsAffected %v", t.String(), t.RowsAffected(), expectedCount)), err0)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err)
	}
	return CommandTag{Sql: t.String(), RowsAffected: t.RowsAffected(), Insert: t.Insert(), Update: t.Update(), Delete: t.Delete(), Select: t.Select()}, runtime.NewStatusOK()
}
