package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
)

var (
	pingLoc = pkgPath + "/stat"
)

func Ping[E template.ErrorHandler](ctx context.Context) (status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	fn, ctx, limited = actuatorApply(ctx, &status, PingUri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return template.NewStatusCode(template.StatusRateLimited)
	}
	if dbClient == nil {
		return e.HandleWithContext(ctx, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	return e.HandleWithContext(ctx, pingLoc, dbClient.Ping(ctx))
}
