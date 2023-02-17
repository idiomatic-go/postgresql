package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/runtime"
	"github.com/idiomatic-go/motif/template"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	statLoc = pkgPath + "/stat"
)

// Stat - templated function for retrieving runtime stats
func Stat[E template.ErrorHandler](ctx context.Context) (stat *pgxpool.Stat, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	fn, ctx, limited = actuatorApply(ctx, messaging.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}
