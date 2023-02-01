package sqldml

import (
	"fmt"
	"strings"
)

const (
	InsertEntryStmt = "INSERT INTO test_entry (id,customer_id,category,traffic_type,traffic_protocol,threshold_percent,name,application,route_name,filter_status_codes,status_codes) VALUES"
	NextValFn       = "nextval('test_entry_Id')"
)

func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

func ExampleWriteInsert() {
	stmt, err := WriteInsert(InsertEntryStmt, []any{100, "test string", false, Function(NextValFn), Function(ChangedTimestampFn)})
	fmt.Printf("Stmt    : %v\n", stmt)
	fmt.Printf("Error   : %v\n", err)

	//Output:
	//Stmt    : INSERT INTO test_entry (id,customer_id,category,traffic_type,traffic_protocol,threshold_percent,name,application,route_name,filter_status_codes,status_codes) VALUES
	//(100,'test string',false,nextval('test_entry_Id'),now());
	//
	//Error   : <nil>

}

func ExampleWriteInsertValues() {
	sb := strings.Builder{}

	err := WriteInsertValues(&sb, nil)
	fmt.Printf("Stmt    : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error   : %v\n", err)

	sb1 := strings.Builder{}
	err = WriteInsertValues(&sb1, []any{100})
	fmt.Printf("Stmt    : %v\n", sb1.String())
	fmt.Printf("Error   : %v\n", err)

	err = WriteInsertValues(&sb, []any{100, "test string", false, Function(NextValFn), Function(ChangedTimestampFn)})
	fmt.Printf("Stmt    : %v\n", sb.String())
	fmt.Printf("Error   : %v\n", err)

	//Output:
	//Stmt    : <nil>
	//Error   : invalid insert argument, values slice is empty
	//Stmt    : (100)
	//Error   : <nil>
	//Stmt    : (100,'test string',false,nextval('test_entry_Id'),now())
	//Error   : <nil>

}
