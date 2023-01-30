package pgxsql

const (
	urnPrefix = "urn:postgres:"
)

type Request struct {
	Uri string
	Sql string
}

func (r Request) String() string {
	return r.Sql
}

func NewRequest(path string, sql string) Request {
	return Request{Uri: BuildUri(path), Sql: sql}
}

func BuildUri(path string) string {
	return urnPrefix + path
}
