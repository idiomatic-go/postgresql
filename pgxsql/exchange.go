package pgxsql

import "context"

// ExecExchange - interface for exec
type ExecExchange interface {
	Exec(req *Request) (CommandTag, error)
}

type execExchange struct {
	exec func(*Request) (CommandTag, error)
}

// NewExecExchange - create a ExecExchange interface
func NewExecExchange(fn func(req *Request) (tag CommandTag, err error)) ExecExchange {
	return &execExchange{fn}
}

// Exec - call the Exec function
func (e *execExchange) Exec(req *Request) (CommandTag, error) {
	return e.exec(req)
}

func execExchangeCast(ctx context.Context) (ExecExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(ExecExchange); ok {
		return e, true
	}
	return nil, false
}

// QueryExchange - interface for query
type QueryExchange interface {
	Query(req *Request) (Rows, error)
}

type queryExchange struct {
	query func(*Request) (Rows, error)
}

// NewQueryExchange - create a QueryExchange interface
func NewQueryExchange(fn func(req *Request) (Rows, error)) QueryExchange {
	return &queryExchange{query: fn}
}

// Query - call the Query function
func (e *queryExchange) Query(req *Request) (Rows, error) {
	return e.query(req)
}

func queryExchangeCast(ctx context.Context) (QueryExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(QueryExchange); ok {
		return e, true
	}
	return nil, false
}
