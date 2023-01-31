package pgxsql

import "fmt"

func ExampleBuildRequest() {
	rsc := "exec-test-resource.dev"
	uri := BuildExecUri(rsc)

	fmt.Printf("test: BuildExecUri(%v) -> %v\n", rsc, uri)

	rsc = "query-test-resource.prod"
	uri = BuildQueryUri(rsc)

	fmt.Printf("test: BuildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: BuildExecUri(exec-test-resource.dev) -> urn:postgres:exec.exec-test-resource.dev
	//test: BuildQueryUri(query-test-resource.prod) -> urn:postgres:query.query-test-resource.prod

}
