-- Example DDL schema for quickgen
-- Run: go run ./cmd/quickgen --schema example.sql

CREATE TABLE IF NOT EXISTS user_profiles (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL UNIQUE,
    full_name   TEXT     NOT NULL,
    avatar_url  TEXT,
    bio         TEXT,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS roles (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    name        TEXT     NOT NULL UNIQUE,
    slug        TEXT     NOT NULL UNIQUE,
    description TEXT,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);
