package db

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

func Connect() (*bun.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "flags.db"
	}

	// WAL mode allows concurrent reads without blocking writes
	sqldb, err := sql.Open("sqlite", dbPath+"?_journal=WAL&_timeout=5000&_fk=true")
	if err != nil {
		return nil, err
	}

	// SQLite only has one writer at a time; cap pool so we don't queue up
	// hundreds of connections all competing for the write lock.
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)

	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}
