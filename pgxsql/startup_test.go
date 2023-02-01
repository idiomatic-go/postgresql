package pgxsql

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/messaging"
	"time"
)

const (
	serviceUrlFmt = "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"
	serviceUrl    = ""
)

func Example_Startup() {
	fmt.Printf("test: IsStarted() -> %v\n", IsStarted())
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())

	}

	//Output:
	//test: IsStarted() -> false
	//test: clientStartup() -> [started:true]

}

func testStartup() error {
	if serviceUrl == "" {
		return errors.New("error running testStartup(): service url is empty")
	}
	c <- messaging.Message{
		To:      "",
		From:    "",
		Event:   messaging.StartupEvent,
		Status:  nil,
		Content: []any{messaging.DatabaseUrl{Url: serviceUrl}},
		ReplyTo: nil,
	}
	time.Sleep(time.Second * 3)
	return nil
}
