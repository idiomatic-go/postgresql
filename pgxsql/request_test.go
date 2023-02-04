package pgxsql

import "fmt"

func ExampleBuildRequest() {
	rsc := "exec-test-resource.dev"
	uri := BuildInsertUri(rsc)

	fmt.Printf("test: BuildInsertUri(%v) -> %v\n", rsc, uri)

	rsc = "query-test-resource.prod"
	uri = BuildQueryUri(rsc)

	fmt.Printf("test: BuildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: BuildInsertUri(exec-test-resource.dev) -> urn:postgres:insert.exec-test-resource.dev
	//test: BuildQueryUri(query-test-resource.prod) -> urn:postgres:query.query-test-resource.prod

}

func ExampleRequest_Validate() {
	uri := "urn:postgres:query.resource"
	sql := "select * from table"
	req := Request{}

	err := req.Validate()
	fmt.Printf("test: Validate(empty) -> %v\n", err)

	req.Uri = uri
	err = req.Validate()
	fmt.Printf("test: Validate(%v) -> %v\n", uri, err)

	req.Uri = ""
	req.Template = sql
	err = req.Validate()
	fmt.Printf("test: Validate(%v) -> %v\n", sql, err)

	req.Uri = uri
	req.Template = sql
	err = req.Validate()
	fmt.Printf("test: Validate(all) -> %v\n", err)

	//Output:
	//test: Validate(empty) -> invalid argument: request Uri is empty
	//test: Validate(urn:postgres:query.resource) -> invalid argument: request template is empty
	//test: Validate(select * from table) -> invalid argument: request Uri is empty
	//test: Validate(all) -> <nil>

}

func ExampleBuildSql() {
	rsc := "access-log"
	t := "delete from access_log"
	req := NewDeleteRequest(rsc, t, nil)

	sql := req.BuildSql()
	fmt.Printf("test: Delete.BuildSql(%v) -> %v\n", t, sql)

	t = "update access_log"
	req = NewUpdateRequest(rsc, t, nil, nil)
	sql = req.BuildSql()
	fmt.Printf("test: Update.BuildSql(%v) -> %v\n", t, sql)

	//Output:
	//test: Delete.BuildSql(delete from access_log) -> delete from access_log
	//WHERE delete_error_no_where_clause = 'null';
	//test: Update.BuildSql(update access_log) -> update access_log
	//SET update_error_no_set_clause = 'null'
	//WHERE update_error_no_where_clause = 'null';
	
}
