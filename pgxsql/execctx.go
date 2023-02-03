package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"time"
)

var (
	execContextKey  = &contextKey{"pgxsql-exec"}
	emptyCommandTag = CommandTag{}
)

type ExecProxy func(req Request) (CommandTag, error)

// ContextWithExec - creates a new Context with an Exec function
func ContextWithExec(ctx context.Context, fn ExecProxy) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &execCtx{ctx, execContextKey, fn}
}

func ContextExec(ctx context.Context, req Request) (CommandTag, error) {
	if ctx == nil {
		return emptyCommandTag, errors.New("context is nil")
	}
	i := ctx.Value(queryContextKey)
	if template.IsNil(i) {
		return emptyCommandTag, errors.New("context value nil")
	}
	if exec, ok := i.(ExecProxy); ok && exec != nil {
		return exec(req)
	}
	return CommandTag{}, errors.New("context value is not of type ExecProxy")
}

func IsContextExec(c context.Context) bool {
	if c == nil {
		return false
	}
	for {
		switch c.(type) {
		case *execCtx:
			return true
		default:
			return false
		}
	}
}

type execCtx struct {
	ctx      context.Context
	key, val any
}

func (*execCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*execCtx) Done() <-chan struct{} {
	return nil
}

func (*execCtx) Err() error {
	return nil
}

func (v *execCtx) Value(key any) any {
	return v.val
}
