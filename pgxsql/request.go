package pgxsql

import (
	"errors"
	"github.com/idiomatic-go/postgresql/pgxdml"
	"net/url"
)

const (
	UrnNID    = "urn:postgres"
	QueryNSS  = UrnNID + ":" + "query."
	InsertNSS = UrnNID + ":" + "insert."
	UpdateNSS = UrnNID + ":" + "update."
	DeleteNSS = UrnNID + ":" + "delete."

	PingUri  = UrnNID + ":" + PingPath
	PingPath = "ping"
	StatUri  = UrnNID + ":" + "stat"

	selectCmd = 0
	insertCmd = 1
	updateCmd = 2
	deleteCmd = 3

	variableReference = "$1"
)

func BuildQueryUri(resource string) string {
	return QueryNSS + resource
}

func BuildInsertUri(resource string) string {
	return InsertNSS + resource
}

func BuildUpdateUri(resource string) string {
	return UpdateNSS + resource
}

func BuildDeleteUri(resource string) string {
	return DeleteNSS + resource
}

type Request struct {
	cmd      int
	Uri      string
	Template string
	Values   [][]any
	Attrs    []pgxdml.Attr
	Where    []pgxdml.Attr
	Error    error
}

func (r *Request) Validate() error {
	if r.Uri == "" {
		return errors.New("invalid argument: request Uri is empty")
	}
	if r.Template == "" {
		return errors.New("invalid argument: request template is empty")
	}
	if r.cmd == deleteCmd {
		if len(r.Where) == 0 {
			return errors.New("invalid argument: delete where clause is empty")
		}
	}
	if r.cmd == updateCmd {
		if len(r.Attrs) == 0 {
			return errors.New("invalid argument: update set clause is empty")
		}
		if len(r.Where) == 0 {
			return errors.New("invalid argument: update where clause is empty")
		}
	}
	return nil
}

func (r *Request) String() string {
	return r.Template
}

func (r *Request) BuildSql() string {
	var sql string
	var err error

	switch r.cmd {
	case selectCmd:
		sql, err = pgxdml.ExpandSelect(r.Template, r.Where)
	case insertCmd:
		sql, err = pgxdml.WriteInsert(r.Template, r.Values)
	case updateCmd:
		if len(r.Where) == 0 {
			r.Where = append(r.Where, pgxdml.Attr{Name: "update_error_no_where_clause", Val: "null"})
		}
		if len(r.Attrs) == 0 {
			r.Attrs = append(r.Attrs, pgxdml.Attr{Name: "update_error_no_set_clause", Val: "null"})
		}
		sql, err = pgxdml.WriteUpdate(r.Template, r.Attrs, r.Where)
	case deleteCmd:
		if len(r.Where) == 0 {
			r.Where = append(r.Where, pgxdml.Attr{Name: "delete_error_no_where_clause", Val: "null"})
		}
		sql, err = pgxdml.WriteDelete(r.Template, r.Where)
	}
	r.Error = err
	return sql
}

func NewQueryRequest(resource, template string, where []pgxdml.Attr) *Request {
	return &Request{cmd: selectCmd, Uri: BuildQueryUri(resource), Template: template, Where: where}
}

func NewQueryRequestFromUrl(resource, template string, url *url.URL) *Request {
	return &Request{cmd: selectCmd, Uri: BuildQueryUri(resource), Template: template, Where: pgxdml.BuildWhere(url)}
}

func NewInsertRequest(resource, template string, values [][]any) *Request {
	return &Request{cmd: insertCmd, Uri: BuildInsertUri(resource), Template: template, Values: values}
}

func NewUpdateRequest(resource, template string, attrs []pgxdml.Attr, where []pgxdml.Attr) *Request {
	return &Request{cmd: updateCmd, Uri: BuildUpdateUri(resource), Template: template, Attrs: attrs, Where: where}
}

func NewDeleteRequest(resource, template string, where []pgxdml.Attr) *Request {
	return &Request{cmd: deleteCmd, Uri: BuildDeleteUri(resource), Template: template, Attrs: nil, Where: where}
}

/*
func NewExecRequest(resource, sql string) Request {
	return Request{Uri: BuildExecUri(resource), Sql: sql}
}

func BuildExecUri(resource string) string {
	return ExecNSS + resource
}


*/
