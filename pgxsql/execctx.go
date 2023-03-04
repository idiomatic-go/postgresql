package pgxsql

import (
	"context"
	"time"
)

/*
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

*/

type ExecContext interface {
	context.Context
	ExecExchange
	withValue(key, val any) context.Context
}

type execContext struct {
	ctx      context.Context
	exchange ExecExchange
}

func NewExecContext(ctx context.Context, exec ExecExchange) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &execContext{ctx: ctx, exchange: exec}
}

func (c *execContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *execContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *execContext) Err() error {
	return c.ctx.Err()
}

func (c *execContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *execContext) Exec(req *Request) (CommandTag, error) {
	return c.exchange.Exec(req)
}

func (c *execContext) withValue(key, val any) context.Context {
	c.ctx = context.WithValue(c.ctx, key, val)
	return c
}

func ExecContextWithValue(ctx context.Context, key any, val any) context.Context {
	if ctx == nil {
		return nil
	}
	if curr, ok := any(ctx).(ExecContext); ok {
		return curr.withValue(key, val)
	}
	return ctx
}

func IsExecContext(ctx context.Context) bool {
	if _, ok := any(ctx).(ExecContext); ok {
		return true
	}
	return false
}
