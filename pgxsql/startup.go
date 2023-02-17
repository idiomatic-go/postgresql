package pgxsql

import (
	"context"
	"github.com/idiomatic-go/motif/messaging"
	"github.com/idiomatic-go/motif/template"
	"reflect"
	"sync/atomic"
	"time"
)

type pkg struct{}

var (
	Uri           = pkgPath
	c             = make(chan messaging.Message, 1)
	pkgPath       = reflect.TypeOf(any(pkg{})).PkgPath()
	started       int64
	actuatorApply messaging.ActuatorApply
)

// IsStarted - returns status of startup
func IsStarted() bool {
	return atomic.LoadInt64(&started) != 0
}

func setStarted() {
	atomic.StoreInt64(&started, 1)
}

func resetStarted() {
	atomic.StoreInt64(&started, 0)
}

func init() {
	actuatorApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
		return func() {}, ctx, false
	}
	messaging.RegisterResource(Uri, c)
	go receive()
}

var messageHandler messaging.MessageHandler = func(msg messaging.Message) {
	switch msg.Event {
	case messaging.StartupEvent:
		clientStartup(msg)
		if IsStarted() {
			apply := messaging.AccessActuatorApply(&msg)
			if apply != nil {
				actuatorApply = apply
			}
		}
	case messaging.ShutdownEvent:
		ClientShutdown()
	case messaging.PingEvent:
		start := time.Now()
		messaging.ReplyTo(msg, Ping[template.LogError](nil).SetDuration(time.Since(start)))
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			go messageHandler(msg)
		default:
		}
	}
}
