package pgxsql

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// psql -x "postgres://tsdbadmin@t9aggksc24.gspnhi29bv.tsdb.cloud.timescale.com:33251/tsdb?sslmode=require"
// Password for user tsdbadmin:

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/messaging"
	"github.com/idiomatic-go/middleware/runtime"
	"github.com/idiomatic-go/middleware/template"
	"github.com/idiomatic-go/postgresql/resource"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	configNameFmt  = "fs/pgxsql/config_%v.json"
	DatabaseURLKey = "DATABASE_URL"
)

var dbClient *pgxpool.Pool
var clientLoc = pkgPath + "/client"

var clientStartup messaging.MessageHandler = func(msg messaging.Message) {
	if IsStarted() {
		return
	}
	start := time.Now()
	name := runtime.EnvExpansion(configNameFmt)
	m, err := resource.ReadMap(name)
	if err != nil {
		messaging.ReplyTo(msg, template.NewStatusError(clientLoc, errors.New(fmt.Sprintf("error reading configuration file: %v", name))).SetDuration(time.Since(start)))
		return
	}
	credentials := messaging.AccessCredentials(&msg)
	if credentials == nil {
		messaging.ReplyTo(msg, template.NewStatusError(clientLoc, errors.New("credentials function is nil")).SetDuration(time.Since(start)))
		return
	}
	err = ClientStartup(m, credentials)
	if err != nil {
		messaging.ReplyTo(msg, template.NewStatusError(clientLoc, err).SetDuration(time.Since(start)))
		return
	}
	messaging.ReplyTo(msg, template.NewStatusOK().SetDuration(time.Since(start)))
}

func ClientStartup(config map[string]string, credentials messaging.Credentials) error {
	if IsStarted() {
		return nil
	}
	// Access database URL
	url, ok := config[DatabaseURLKey]
	if !ok || url == "" {
		return errors.New(fmt.Sprintf("database URL does not exist in map, or value is empty : %v\n", DatabaseURLKey))
	}
	// Create connection string with credentials
	s, err := connectString(url, credentials)
	if err != nil {
		return err
	}
	// Create pooled client and acquire connection
	dbClient, err = pgxpool.New(context.Background(), s)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create connection pool: %v\n", err))
	}
	conn, err1 := dbClient.Acquire(context.Background())
	if err1 != nil {
		ClientShutdown()
		return errors.New(fmt.Sprintf("unable to acquire connection from pool: %v\n", err1))
	}
	conn.Release()
	setStarted()
	return nil
}

func ClientShutdown() {
	if dbClient != nil {
		resetStarted()
		dbClient.Close()
		dbClient = nil
	}
}

func connectString(url string, credentials messaging.Credentials) (string, error) {
	// Username and password can b in the connect string Url
	if credentials == nil {
		return url, nil
	}
	username, password, err := credentials()
	if err != nil {
		return "", errors.New(fmt.Sprintf("error accessing credentials: %v\n", err))
	}
	return fmt.Sprintf(url, username, password), nil
}
