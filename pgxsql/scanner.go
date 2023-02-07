package pgxsql

import (
	"errors"
)

type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
}

func Scan[T Scanner[T]](rows Rows) ([]T, error) {
	if rows == nil {
		return nil, errors.New("invalid request: rows interface is nil")
	}
	var s T
	var t []T
	var err error
	var values []any

	defer rows.Close()
	names := createColumnNames(rows.FieldDescriptions())
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return t, err
		}
		values, err = rows.Values()
		if err != nil {
			return t, err
		}
		convertTimestamps(rows.FieldDescriptions(), values)
		val, err1 := s.Scan(names, values)
		if err1 != nil {
			return t, err1
		}
		t = append(t, val)
	}
	return t, nil
}

func createColumnNames(fields []FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}

func convertTimestamps(fields []FieldDescription, values []any) {
	for i, fld := range fields {
		v := values[i]
		if v != nil {
		}
		if fld.DataTypeOID == 1184 {

		}
	}

}
