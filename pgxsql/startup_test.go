package pgxsql

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/messaging"
	"github.com/idiomatic-go/middleware/template"
	"github.com/idiomatic-go/resiliency/actuator"
	"time"
)

// "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"

const (
	serviceUrl = "postgres://tsdbadmin:akmoebwl97kl31yl@xz90jtzdq1.q2h8vd0pwk.tsdb.cloud.timescale.com:31770/tsdb?sslmode=require"
)

func Example_Startup() {
	//fmt.Printf("test: Uri -> %v\n", Uri)
	fmt.Printf("test: IsStarted() -> %v\n", IsStarted())
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())

		status := Ping[template.DebugError](nil)
		fmt.Printf("test: messaging.Ping() -> %v\n", status)

	}

	//Output:
	//test: IsStarted() -> false
	//test: clientStartup() -> [started:true]
	//{traffic:egress, route:*, request-id:, status-code:0, method:GET, url:urn:postgres:ping, host:postgres, path:ping, timeout:-1, rate-limit:-1, rate-burst:-1, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:}
	//test: messaging.Ping() -> OK

}

func testStartup() error {
	if serviceUrl == "" {
		return errors.New("error running testStartup(): service url is empty")
	}
	if IsStarted() {
		return nil
	}
	c <- messaging.Message{
		To:      "",
		From:    "",
		Event:   messaging.StartupEvent,
		Status:  nil,
		Content: []any{messaging.DatabaseUrl{Url: serviceUrl}, messaging.ActuatorApply(actuator.EgressApply)},
		ReplyTo: nil,
	}
	time.Sleep(time.Second * 3)
	return nil
}
