package pgxsql

type CollectableRow interface {
	FieldDescriptions() []FieldDescription
	Scan(dest ...any) error
	Values() ([]any, error)
	RawValues() [][]byte
}

// RowToFunc is a function that scans or otherwise converts row to a T.
type RowToFunc[T any] func(row CollectableRow) (T, error)

// CollectRows iterates through rows, calling fn for each row, and collecting the results into a slice of T.
func CollectRows[T any](rows Rows, fn RowToFunc[T]) ([]T, error) {
	defer rows.Close()

	slice := []T{}

	for rows.Next() {
		value, err := fn(rows)
		if err != nil {
			return nil, err
		}
		slice = append(slice, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slice, nil
}
