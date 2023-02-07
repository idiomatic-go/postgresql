package pgxsql

import (
	"fmt"
	"github.com/idiomatic-go/motif/messaging"
)

func ExampleClientStartup() {

	db := messaging.DatabaseUrl{Url: ""}
	err := ClientStartup(db, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
