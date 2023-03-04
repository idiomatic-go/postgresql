package pgxsql

import (
	"context"
	"time"
)

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
