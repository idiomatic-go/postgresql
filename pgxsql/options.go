package pgxsql

import "github.com/idiomatic-go/middleware/template"

func SetActuator(fn template.ActuatorApply) {
	if fn != nil {
		actuatorApply = fn
	}
}
