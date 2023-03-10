package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/runtime"
)

var (
	statLoc = pkgPath + "/stat"
)

// Stat - templated function for retrieving runtime stats
func Stat[E runtime.ErrorHandler](ctx context.Context) (stat *Stats, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	fn, ctx, limited = controllerApply(ctx, messaging.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}
