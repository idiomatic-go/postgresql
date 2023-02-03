package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"github.com/jackc/pgx/v5/pgconn"
)

func Query[E template.ErrorHandler](ctx context.Context, req Request, args ...any) (result Rows, status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = actuatorApply(ctx, &status, req.Uri, template.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, template.NewStatusCode(template.StatusRateLimited)
	}
	if IsContextQuery(ctx) {
		var err error
		result, err = ContextQuery(ctx, req)
		return result, e.HandleWithContext(ctx, execLoc, err)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetCode(template.StatusInvalidArgument)
	}
	pgxRows, err := dbClient.Query(ctx, req.Sql, args...)
	if err != nil {
		return nil, e.HandleWithContext(ctx, queryLoc, recast(err))
	}
	return &rows{pgxRows: pgxRows, fd: fieldDescriptions(pgxRows.FieldDescriptions())}, template.NewStatusOK()
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
