package pgxsql

import (
	"context"
	"errors"
	"github.com/idiomatic-go/middleware/template"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FieldDescription struct {
	Name                 string
	TableOID             uint32
	TableAttributeNumber uint16
	DataTypeOID          uint32
	DataTypeSize         int16
	TypeModifier         int32
	Format               int16
}

type Rows interface {
	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() error

	// CommandTag returns the command tag from this query. It is only available after Rows is closed.
	CommandTag() CommandTag

	FieldDescriptions() []FieldDescription

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	// Scan reads the values from the current row into dest values positionally.
	// dest can include pointers to core sqldml, values implementing the Scanner
	// interface, and nil. nil will skip the value entirely. It is an error to
	// call Scan without first calling Next() and checking that it returned true.
	Scan(dest ...any) error

	// Values returns the decoded row values. As with Scan(), it is an error to
	// call Values without first calling Next() and checking that it returned
	// true.
	Values() ([]any, error)

	// RawValues returns the unparsed bytes of the row values. The returned data is only valid until the next Next
	// call or the Rows is closed.
	RawValues() [][]byte

	// Conn returns the underlying *Conn on which the queryv1 was executed. This may return nil if Rows did not come from a
	// *Conn (e.g. if it was created by RowsFromResultReader)
	// TODO : determine use case
	//Conn() *Conn
}

var (
	queryLoc = pkgPath + "/query"
)

type rows struct {
	pgxRows pgx.Rows
	fd      []FieldDescription
}

func (r *rows) Close() {
	if r != nil {
		r.pgxRows.Close()
	}
}

func (r *rows) Err() error {
	if r == nil {
		return nil
	}
	return r.pgxRows.Err()
}

func (r *rows) CommandTag() CommandTag {
	if r == nil {
		return CommandTag{}
	}
	t := r.pgxRows.CommandTag()
	return CommandTag{RowsAffected: t.RowsAffected(), Sql: t.String()}
}

func (r *rows) FieldDescriptions() []FieldDescription {
	if r == nil {
		return nil
	}
	return r.fd
}

func (r *rows) Next() bool {
	if r == nil {
		return false
	}
	return r.pgxRows.Next()
}

func (r *rows) Scan(dest ...any) error {
	if r == nil {
		return nil
	}
	return r.pgxRows.Scan(dest)
}

func (r *rows) Values() ([]any, error) {
	if r == nil {
		return nil, nil
	}
	return r.pgxRows.Values()
}

func (r *rows) RawValues() [][]byte {
	if r == nil {
		return nil
	}
	return r.pgxRows.RawValues()
}

func Query[E template.ErrorHandler](ctx context.Context, req Request, args ...any) (result Rows, status *template.Status) {
	var e E
	var limited = false
	var fn template.ActuatorComplete

	ctx = template.IfElse[context.Context](ctx != nil, ctx, context.Background())
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
		return nil, e.HandleWithContext(ctx, queryLoc, err)
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
