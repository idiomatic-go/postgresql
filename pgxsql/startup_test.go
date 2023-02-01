package pgxsql

import (
	"fmt"
	"github.com/idiomatic-go/middleware/messaging"
	"time"
)

const (
	serviceUrl = "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"
)

func Example_Startup() {
	fmt.Printf("test: IsStarted() -> %v\n", IsStarted())
	startup()
	fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())
	defer ClientShutdown()

	//Output:
	//test: IsStarted() -> false
	//test: clientStartup() -> [started:true]

}

func startup() {
	c <- messaging.Message{
		To:      "",
		From:    "",
		Event:   messaging.StartupEvent,
		Status:  nil,
		Content: []any{messaging.DatabaseUrl{Url: serviceUrl}},
		ReplyTo: nil,
	}
	time.Sleep(time.Second * 3)
}
