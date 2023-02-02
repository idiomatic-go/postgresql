package pgxsql

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func recast(err error) error {
	if err == nil {
		return err
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		msg := fmt.Sprintf("code:%v, message:%v, position:%v, line:%v",
			pgErr.Code, pgErr.Message, pgErr.InternalPosition, pgErr.Line)
		//pgErr.Code,pgErr.ColumnName,pgErr.ConstraintName,pgErr.DataTypeName,pgErr.Message,
		//pgErr.Detail,pgErr.File,pgErr.Hint,pgErr.InternalPosition,pgErr.Line
		return errors.New(msg)
	}
	return err
}
