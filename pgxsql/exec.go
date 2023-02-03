package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/template"
	"github.com/idiomatic-go/postgresql/pgxdml"
)

type CommandTag struct {
	Sql          string
	RowsAffected int64
	Insert       bool
	Update       bool
	Delete       bool
	Select       bool
}

const (
	NullCount = int64(-1)
)

var execLoc = pkgPath + "/exec"

func Insert[E template.ErrorHandler](ctx context.Context, expectedCount int64, req Request, values [][]any) (CommandTag, *template.Status) {
	var e E

	if IsContextExec(ctx) {
		newTag, err := ContextExec(ctx, req)
		return newTag, e.HandleWithContext(ctx, execLoc, err)
	}
	stmt, err := pgxdml.WriteInsert(req.Sql, values)
	if err != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err)
	}
	return Exec[E](ctx, expectedCount, Request{Uri: req.Uri, Sql: stmt})
}

func Update[E template.ErrorHandler](ctx context.Context, expectedCount int64, req Request, attrs []pgxdml.Attr, where []pgxdml.Attr) (CommandTag, *template.Status) {
	var e E

	if IsContextExec(ctx) {
		newTag, err := ContextExec(ctx, req)
		return newTag, e.HandleWithContext(ctx, execLoc, err)
	}
	stmt, err := pgxdml.WriteUpdate(req.Sql, attrs, where)
	if err != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err)
	}
	return Exec[E](ctx, expectedCount, Request{Uri: req.Uri, Sql: stmt})
}

func Delete[E template.ErrorHandler](ctx context.Context, expectedCount int64, req Request, where []pgxdml.Attr) (CommandTag, *template.Status) {
	var e E

	if IsContextExec(ctx) {
		newTag, err := ContextExec(ctx, req)
		return newTag, e.HandleWithContext(ctx, execLoc, err)
	}
	stmt, err := pgxdml.WriteDelete(req.Sql, where)
	if err != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err)
	}
	return Exec[E](ctx, expectedCount, Request{Uri: req.Uri, Sql: stmt})
}

func Exec[E template.ErrorHandler](ctx context.Context, expectedCount int64, req Request, args ...any) (_ CommandTag, status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = actuatorApply(ctx, &status, req.Uri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return CommandTag{}, template.NewStatusCode(template.StatusRateLimited)
	}
	if IsContextExec(ctx) {
		result, err := ContextExec(ctx, req)
		return result, e.HandleWithContext(ctx, execLoc, err)
	}
	if dbClient == nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return CommandTag{}, e.HandleWithContext(ctx, execLoc, err0)
	}
	t, err := dbClient.Exec(ctx, req.Sql, args...)
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
	return CommandTag{Sql: t.String(), RowsAffected: t.RowsAffected(), Insert: t.Insert(), Update: t.Update(), Delete: t.Delete(), Select: t.Select()}, template.NewStatusOK()
}
