package db

import (
	"database/sql"
	"fmt"
)

// Migrate runs all startup migrations. Safe to call on every server start.
func Migrate(db *sql.DB) error {
	if err := migrateUsers(db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func migrateUsers(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			email      TEXT    NOT NULL UNIQUE,
			name       TEXT    NOT NULL,
			avatar_url TEXT,
			provider   TEXT    NOT NULL DEFAULT 'google',
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
		)
	`)
	if err != nil {
		return fmt.Errorf("users table: %w", err)
	}
	return nil
}
