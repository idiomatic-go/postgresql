package sqldml

import (
	"fmt"
	"strings"
)

const (
	UpdateTestEntryStmt = "UPDATE test_entry"
)

func ExampleWriteUpdate() {
	where := []Attr{{Name: "customer_id", Val: "customer1"}, {Name: "created_ts", Val: "2022/11/30 15:48:54.049496"}} //time.Now()}}
	attrs := []Attr{{Name: "status_code", Val: "503"}, {Name: "minimum_code", Val: 99}, {Name: "created_ts", Val: Function("now()")}}

	sql, err := WriteUpdate(UpdateTestEntryStmt, attrs, where)
	fmt.Printf("Stmt       : %v\n", NilEmpty(sql))
	fmt.Printf("Error      : %v\n", err)

	//Output:
	//Stmt       : UPDATE test_entry
	//SET status_code = '503',
	//minimum_code = 99,
	//created_ts = now()
	//WHERE customer_id = 'customer1' AND created_ts = '2022/11/30 15:48:54.049496';
	//Error      : <nil>
}

func ExampleWriteUpdateSet() {
	sb := strings.Builder{}

	err := WriteUpdateSet(&sb, nil)
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	sb.Reset()
	err = WriteUpdateSet(&sb, []Attr{{Name: "status_code", Val: "503"}})
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	sb.Reset()
	err = WriteUpdateSet(&sb, []Attr{{Name: "status_code", Val: "503"}, {Name: "minimum_code", Val: 99}, {Name: "created_ts", Val: Function("now()")}})
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	//Output:
	//Stmt       : <nil>
	//Error      : invalid update set argument, attrs slice is empty
	//Stmt       : SET status_code = '503'
	//
	//Error      : <nil>
	//Stmt       : SET status_code = '503',
	//minimum_code = 99,
	//created_ts = now()
	//
	//Error      : <nil>
}

func ExampleWriteUpdateWhere() {
	sb := strings.Builder{}

	err := WriteUpdateWhere(&sb, nil)
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	err = WriteUpdateWhere(&sb, []Attr{{Name: "", Val: nil}})
	fmt.Printf("Stmt       : %v\n", NilEmpty(strings.Trim(sb.String(), " ")))
	fmt.Printf("Error      : %v\n", err)

	sb.Reset()
	err = WriteUpdateWhere(&sb, []Attr{{Name: "status_code", Val: "503"}})
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	sb.Reset()
	err = WriteUpdateWhere(&sb, []Attr{{Name: "status_code", Val: "503"}, {Name: "minimum_code", Val: 99}, {Name: "created_ts", Val: Function("now()")}})
	fmt.Printf("Stmt       : %v\n", NilEmpty(sb.String()))
	fmt.Printf("Error      : %v\n", err)

	//Output:
	//Stmt       : <nil>
	//Error      : invalid update where argument, attrs slice is empty
	//Stmt       : WHERE
	//Error      : invalid attribute argument, attribute name is empty
	//Stmt       : WHERE status_code = '503';
	//Error      : <nil>
	//Stmt       : WHERE status_code = '503' AND minimum_code = 99 AND created_ts = now();
	//Error      : <nil>

}
