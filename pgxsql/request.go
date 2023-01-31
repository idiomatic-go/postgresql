package pgxsql

const (
	UrnNID   = "urn:postgres"
	QueryNSS = UrnNID + ":" + "query."
	ExecNSS  = UrnNID + ":" + "exec."
	PingUri  = UrnNID + ":" + "ping"
	StatUri  = UrnNID + ":" + "stat"
)

type Request struct {
	Uri string
	Sql string
}

func (r Request) String() string {
	return r.Sql
}

func NewQueryRequest(resource, sql string) Request {
	return Request{Uri: BuildQueryUri(resource), Sql: sql}
}

func BuildQueryUri(resource string) string {
	return QueryNSS + resource
}

func NewExecRequest(resource, sql string) Request {
	return Request{Uri: BuildExecUri(resource), Sql: sql}
}

func BuildExecUri(resource string) string {
	return ExecNSS + resource
}
