package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/runtime"
	"github.com/idiomatic-go/motif/template"
	"github.com/jackc/pgx/v5/pgconn"
)

// Query - templated function for a Query
func Query[E template.ErrorHandler](ctx context.Context, req *Request, args ...any) (result Rows, status *runtime.Status) {
	var e E
	var limited = false
	var fn messaging.ActuatorComplete

	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return nil, e.HandleWithContext(ctx, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	fn, ctx, limited = actuatorApply(ctx, &status, req.Uri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if IsContextQuery(ctx) {
		var err error
		result, err = ContextQuery(ctx, req)
		return result, e.HandleWithContext(ctx, execLoc, err)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	pgxRows, err := dbClient.Query(ctx, req.BuildSql(), args...)
	if err != nil {
		return nil, e.HandleWithContext(ctx, queryLoc, recast(err))
	}
	return &proxyRows{pgxRows: pgxRows, fd: fieldDescriptions(pgxRows.FieldDescriptions())}, runtime.NewStatusOK()
}

func fieldDescriptions(fields []pgconn.FieldDescription) []FieldDescription {
	var result []FieldDescription
	for _, f := range fields {
		result = append(result, FieldDescription{Name: f.Name,
			TableOID:             f.TableOID,
			TableAttributeNumber: f.TableAttributeNumber,
			DataTypeOID:          f.DataTypeOID,
			DataTypeSize:         f.DataTypeSize,
			TypeModifier:         f.TypeModifier,
			Format:               f.Format})
	}
	return result
}
