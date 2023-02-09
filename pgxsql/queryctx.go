package pgxsql

import (
	"context"
	"errors"
	"time"
)

var (
	queryContextKey = &contextKey{"pgxsql-query"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

// QueryProxy - proxy type for the Query function
type QueryProxy func(req *Request) (Rows, error)

// ContextWithQuery - creates a new Context with a Query proxy
func ContextWithQuery(ctx context.Context, fn QueryProxy) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if fn == nil {
		return ctx
	}
	return &queryCtx{ctx, queryContextKey, fn}
}

// ContextQuery - calls the Query proxy
func ContextQuery(ctx context.Context, req *Request) (Rows, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	i := ctx.Value(queryContextKey)
	if i == nil {
		return nil, errors.New("context value is nil")
	}
	if query, ok := i.(QueryProxy); ok && query != nil {
		return query(req)
	}
	return nil, errors.New("context value is not of QueryProxy type")
}

// IsContextQuery - determines if the context is a Query proxy context
func IsContextQuery(c context.Context) bool {
	if c == nil {
		return false
	}
	for {
		switch c.(type) {
		case *queryCtx:
			return true
		default:
			return false
		}
	}
}

type queryCtx struct {
	ctx      context.Context
	key, val any
}

func (*queryCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*queryCtx) Done() <-chan struct{} {
	return nil
}

func (*queryCtx) Err() error {
	return nil
}

func (v *queryCtx) Value(key any) any {
	return v.val
}
