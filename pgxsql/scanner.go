package pgxsql

import (
	"errors"
)

type Scanner[T any] interface {
	Scan(fields []FieldDescription, values []any) (T, error)
}

func Scan[T Scanner[T]](rows Rows) ([]T, error) {
	if rows == nil {
		return nil, errors.New("invalid request: rows interface is nil")
	}
	var s T
	var t []T
	var err error
	var values []any

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return t, err
		}
		values, err = rows.Values()
		if err != nil {
			return t, err
		}
		val, err1 := s.Scan(rows.FieldDescriptions(), values)
		if err1 != nil {
			return t, err1
		}
		t = append(t, val)
	}
	return t, nil
}
