package pgxsql

import (
	"fmt"
	"github.com/idiomatic-go/middleware/template"
)

func ExampleStat() {
	startup()
	fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())
	defer ClientShutdown()
	stat, status := Stat[template.DebugError](nil)
	fmt.Printf("test: Stat(nil) -> [status:%v] [stat:%v]\n", status, stat != nil)

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Stat(nil) -> [status:OK] [stat:true]

}
