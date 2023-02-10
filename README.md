# postgresql

## pgxdml

[PostgresDML][pgxdmlpkg] implements types that build SQL statements based on the configured attributes. Support is also available for selecting
PostgreSQL functions for timestamps and next values.

## pgxsql

[PostgresSQL][pgxsqlpkg] provides the templated functions for query, exec, ping, and stat. Testing proxies are implemented for exec and query functions.
The processing of host generated messaging for startup and ping events is also supported. Scanning of PostgreSQL rows into application types utilizes a
templated interface, and corresponding templated Scan function. Care was taken to not leak any direct references to PostgresSQL specific packages.

~~~
// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows Rows) ([]T, error) {
    // implementation details
}
~~~

Resiliency for PostgresSQL database client calls is provided by an [Actuator][actuatorcall] function call that is initialized by the host on startup:
~~~
var actuatorApply messaging.ActuatorApply

fn, ctx, limited = actuatorApply(ctx, &status, req.Uri, runtime.ContextRequestId(ctx), "GET")
defer fn()
~~~

[pgxdmlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxdml/http>
[pgxsqlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxsql>
[actuatorcall]: <https://pkg.go.dev/github.com/idiomatic-go/resiliency/actuator#EgressApply>
