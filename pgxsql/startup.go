package pgxsql

import (
	"context"
	"github.com/idiomatic-go/middleware/messaging"
	"github.com/idiomatic-go/middleware/template"
	"reflect"
	"sync/atomic"
)

type pkg struct{}

var (
	Uri           = pkgPath
	c             = make(chan messaging.Message, 1)
	pkgPath       = reflect.TypeOf(any(pkg{})).PkgPath()
	started       int64
	actuatorApply template.ActuatorApply
)

func IsStarted() bool { return atomic.LoadInt64(&started) != 0 }

func setStarted() {
	atomic.StoreInt64(&started, 1)
}

func resetStarted() {
	atomic.StoreInt64(&started, 0)
}

func complete() {}

func init() {
	actuatorApply = func(ctx context.Context, status **template.Status, uri, requestId, method string) (template.ActuatorComplete, context.Context, bool) {
		return complete, ctx, false
	}
	messaging.RegisterResource(Uri, c)
	go receive()
}

var messageHandler messaging.MessageHandler = func(msg messaging.Message) {
	switch msg.Event {
	case messaging.StartupEvent:
		clientStartup(msg)
	case messaging.ShutdownEvent:
		ClientShutdown()
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			messageHandler(msg)
		default:
		}
	}
}
