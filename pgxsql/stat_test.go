package pgxsql

import (
	"fmt"
	"github.com/idiomatic-go/middleware/template"
)

func ExampleStat() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())

		stat, status := Stat[template.DebugError](nil)
		fmt.Printf("test: Stat(nil) -> [status:%v] [stat:%v]\n", status, stat != nil)
	}

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Stat(nil) -> [status:OK] [stat:true]

}
