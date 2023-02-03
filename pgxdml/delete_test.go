package pgxdml

import "fmt"

const (
	deleteTestEntryStmt = "DELETE test_entry"
)

func ExampleWriteDelete() {
	where := []Attr{{Name: "customer_id", Val: "customer1"}, {Name: "created_ts", Val: "2022/11/30 15:48:54.049496"}} //time.Now()}}

	sql, err := WriteDelete(deleteTestEntryStmt, where)
	fmt.Printf("Stmt       : %v\n", NilEmpty(sql))
	fmt.Printf("Error      : %v\n", err)

	//Output:
	//Stmt       : DELETE test_entry
	//WHERE customer_id = 'customer1' AND created_ts = '2022/11/30 15:48:54.049496';
	//Error      : <nil>

}
