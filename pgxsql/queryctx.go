package pgxsql

import (
	"context"
	"time"
)

type QueryContext interface {
	context.Context
	QueryExchange
	withValue(key, val any) context.Context
}

type queryContext struct {
	ctx      context.Context
	exchange QueryExchange
}

func NewQueryContext(ctx context.Context, query QueryExchange) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &queryContext{ctx: ctx, exchange: query}
}

func (c *queryContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *queryContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *queryContext) Err() error {
	return c.ctx.Err()
}

func (c *queryContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *queryContext) Query(req *Request) (Rows, error) {
	return c.exchange.Query(req)
}

func (c *queryContext) withValue(key, val any) context.Context {
	c.ctx = context.WithValue(c.ctx, key, val)
	return c
}

func QueryContextWithValue(ctx context.Context, key any, val any) context.Context {
	if ctx == nil {
		return nil
	}
	if curr, ok := any(ctx).(QueryContext); ok {
		return curr.withValue(key, val)
	}
	return ctx
}

func IsQueryContext(ctx context.Context) bool {
	if _, ok := any(ctx).(QueryContext); ok {
		return true
	}
	return false
}
