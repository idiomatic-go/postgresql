package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"time"
)

var (
	queryContextKey = &contextKey{"pgxsql-query"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

type QueryProxy func(req *Request) (Rows, error)

// ContextWithQuery - creates a new Context with a Query function
func ContextWithQuery(ctx context.Context, fn QueryProxy) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &queryCtx{ctx, queryContextKey, fn}
}

func ContextQuery(ctx context.Context, req *Request) (Rows, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	i := ctx.Value(queryContextKey)
	if template.IsNil(i) {
		return nil, errors.New("context value is nil")
	}
	if query, ok := i.(QueryProxy); ok && query != nil {
		return query(req)
	}
	return nil, errors.New("context value is not of QueryProxy type")
}

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
