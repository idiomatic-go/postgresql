# postgresql

## pgxdml

[PostgresDML][pgxdmlpkg] implements types that build SQL statements based on the configured attributes. Support is also available for selecting
PostgreSQL functions for timestamps and next values when needed for statement creation.

## pgxsql

[PostgresSQL][pgxsqlpkg] provides the templated functions for query, exec, ping, and stat. Testing proxies are implemented for exec and query functions.
The processing of host generated messaging for startup and ping events is also supported.


[pgxdmlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxdml/http>
[pgxsqlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxsql>
