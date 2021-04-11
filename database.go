package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB DBObj

type DBObj struct {
	Connection *sqlx.DB
}

func initDB(typeDB, conn string) {
	var (
		err error
	)

	DB.Connection, err = sqlx.Connect(typeDB, conn)
	if err != nil {
		log.Fatalf("[InitDB][sqlx.Connect] Input: %v Output: %v", conn, err)
	}

	if _, err = DB.QueryContext(context.Background(), "SELECT 1 FROM user", nil); err != nil {
		log.Fatalf("[InitDB][QueryContext] Input: %v Output: %v", conn, err)
	}
}

func (db *DBObj) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.Connection.QueryRowContext(ctx, query, args...)
}

func (db *DBObj) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Connection.QueryContext(ctx, query, args...)
}

func (db *DBObj) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Connection.ExecContext(ctx, query, args...)
}
