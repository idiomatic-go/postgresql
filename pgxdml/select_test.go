package pgxdml

import "fmt"

func ExampleExpandSelect() {
	t := "select * from access_log where $1 order by start_time desc limit 5"
	where := []Attr{{Name: "status_code", Val: "503"}}

	sql := ExpandSelect("", nil)
	fmt.Printf("test: ExpandSelect(nil,nil) -> [empty:%v]\n", sql == "")

	sql = ExpandSelect(t, nil)
	fmt.Printf("test: ExpandSelect(t,nil) -> %v\n", sql)

	sql = ExpandSelect(t, where)
	fmt.Printf("test: ExpandSelect(t,where) -> %v\n", sql)

	//Output:
	//test: ExpandSelect(nil,nil) -> [empty:true]
	//test: ExpandSelect(t,nil) -> select * from access_log where $1 order by start_time desc limit 5
	//test: ExpandSelect(t,where) -> select * from access_log where status_code = '503' order by start_time desc limit 5

}
