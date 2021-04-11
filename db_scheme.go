package main

const (
	dbTableUser = `
	CREATE TABLE "user" (
		"uid"	TEXT NOT NULL UNIQUE,
		"credential_type"	INTEGER DEFAULT 1,
		"credential_access"	TEXT UNIQUE,
		"credential_key"	TEXT,
		"credential_salt"	TEXT,
		"create_time"	TEXT,
		"update_time"	TEXT,
		PRIMARY KEY("uid")
	);`
)
