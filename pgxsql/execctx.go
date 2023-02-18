package pgxsql

import (
	"context"
	"errors"
	"time"
)

var (
	execContextKey = &contextKey{"pgxsql-exec"}
)

// ExecProxy - proxy type for the Exec function
type ExecProxy func(req *Request) (CommandTag, error)

// ContextWithExec - creates a new Context with an Exec proxy
func ContextWithExec(ctx context.Context, fn ExecProxy) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if fn == nil {
		return ctx
	}
	return &execCtx{ctx, execContextKey, fn}
}

// ContextExec - calls the Exec proxy
func ContextExec(ctx context.Context, req *Request) (tag CommandTag, err error) {
	if ctx == nil {
		return tag, errors.New("context is nil")
	}
	i := ctx.Value(queryContextKey)
	if i == nil {
		return tag, errors.New("context value nil")
	}
	if exec, ok := i.(ExecProxy); ok && exec != nil {
		return exec(req)
	}
	return tag, errors.New("context value is not of type ExecProxy")
}

// IsContextExec - determines if the context is an Exec proxy context
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
