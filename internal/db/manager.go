package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "db/local.db"
	dbName     = "sqlite3"
)

type Manager struct {
	client *sql.DB
	ctx    context.Context
}

func New(ctx context.Context) *Manager {
	client, err := sql.Open(dbName, dbFileName)
	if err != nil {
		panic(fmt.Errorf("cannot create db file %s. error: %v", dbFileName, err))
	}

	if err = client.PingContext(ctx); err != nil {
		panic(fmt.Errorf("cannot connect to db file %s. error: %v", dbFileName, err))
	}

	return &Manager{
		client: client,
		ctx:    ctx,
	}
}
