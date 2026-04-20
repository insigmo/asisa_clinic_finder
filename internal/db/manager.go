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
}

func New(ctx context.Context) (*Manager, error) {
	client, err := sql.Open(dbName, dbFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open db file %s: %w", dbFileName, err)
	}

	if err = client.PingContext(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("cannot connect to db file %s: %w", dbFileName, err)
	}

	return &Manager{client: client}, nil
}

func (db *Manager) Close() error {
	return db.client.Close()
}
