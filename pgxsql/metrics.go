package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	metricLoc = pkgPath + "/stat"
	pingUri   = BuildUri("metrics.ping")
	statUri   = BuildUri("metrics.stat")
)

func Stat[E template.ErrorHandler](ctx context.Context) (stat *pgxpool.Stat, status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	fn, ctx, limited = actuatorApply(ctx, &status, statUri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, template.NewStatusCode(template.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, metricLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	return dbClient.Stat(), nil
}

func Ping[E template.ErrorHandler](ctx context.Context) (status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	fn, ctx, limited = actuatorApply(ctx, &status, pingUri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return template.NewStatusCode(template.StatusRateLimited)
	}
	if dbClient == nil {
		return e.HandleWithContext(ctx, metricLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	return e.HandleWithContext(ctx, metricLoc, dbClient.Ping(ctx))
}
