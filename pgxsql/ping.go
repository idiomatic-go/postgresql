package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/runtime"
	"github.com/idiomatic-go/motif/template"
)

var (
	pingLoc = pkgPath + "/stat"
)

// Ping - templated function for pinging the database cluster
func Ping[E template.ErrorHandler](ctx context.Context) (status *runtime.Status) {
	var e E
	var limited = false
	var fn messaging.ActuatorComplete

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = actuatorApply(ctx, &status, PingUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return e.HandleWithContext(ctx, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return e.HandleWithContext(ctx, pingLoc, dbClient.Ping(ctx))
}
