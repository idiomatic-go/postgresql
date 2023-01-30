package pgxsql

import "fmt"

func Example_Startup() {
	m := map[string]string{DatabaseURLKey: ""}

	err := ClientStartup(m, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL does not exist in map, or value is empty : DATABASE_URL

}
