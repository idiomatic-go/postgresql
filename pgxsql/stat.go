package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	statLoc = pkgPath + "/stat"
)

func Stat[E template.ErrorHandler](ctx context.Context) (stat *pgxpool.Stat, status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	fn, ctx, limited = actuatorApply(ctx, &status, StatUri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, template.NewStatusCode(template.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	return dbClient.Stat(), nil
}
