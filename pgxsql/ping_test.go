package pgxsql

import (
	"fmt"
	"github.com/idiomatic-go/middleware/template"
)

func ExamplePing() {
	startup()
	fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())
	defer ClientShutdown()
	status := Ping[template.DebugError](nil)
	fmt.Printf("test: Ping(nil) -> %v\n", status)

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Ping(nil) -> OK

}
