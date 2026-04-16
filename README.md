# Quick Go

A Go starter project with Auth by Email & OAuth2 via Google. Includes a code generator (`quickgen`) that scaffolds CRUD handlers, repositories, models, and Hugo content pages from a SQL DDL schema.

## Features

- Google OAuth2 login
- JWT-based session with Redis revocation
- SQLite database with auto-migration
- `quickgen` — generates CRUD boilerplate from a SQL schema
- Hugo-based static frontend

## Prerequisites

- Go 1.21+
- Redis (optional — disables session revocation if unavailable)
- Hugo (optional — required to build the frontend)

## Setup

```sh
cp .env.sample .env
# Edit .env and fill in your values
```

## Build

Build both the server and code generator:

```sh
make build
```

Build individually:

```sh
make build-server    # outputs bin/server
make build-quickgen  # outputs bin/quickgen
```

## Run

```sh
make run
```

Or run directly without building:

```sh
make dev
```

The server listens on the port defined in `.env` (default: `8080`).

## Code Generation

Generate CRUD scaffolding from a SQL DDL file:

```sh
make gen SCHEMA=path/to/schema.sql
```

This is equivalent to:

```sh
./bin/quickgen --schema path/to/schema.sql --out .
```

Generated files:

```
generated/
  model/        # struct definitions per table
  repo/         # database query functions per table
  handler/      # HTTP handlers per table
  routes.go     # route registration for all tables
frontend/hugo-site/content/
  <table>/
    _index.md   # list page
    form.md     # form page
```

After generating, register the routes in `cmd/server/main.go`:

```go
generated.RegisterRoutes(r, database)
```

Then rebuild:

```sh
make build-server
```

## Other Commands

```sh
make tidy    # go mod tidy + verify
make clean   # remove bin/
make help    # list available targets
```

## Project Structure

```
cmd/
  server/     # main server entrypoint
  quickgen/   # code generator CLI
internal/
  auth/       # Google OAuth2, JWT, session store
  config/     # env-based config loader
  db/         # SQLite open + migrations
  middleware/ # JWT auth middleware
templates/    # Go text/template files used by quickgen
generated/    # output directory for quickgen (gitignored)
```

## Starting a New Project from This Template

- Set up a Google OAuth2 app and fill in `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` in `.env`
- Write your DB schema as a SQL DDL file
- Run `make gen SCHEMA=your_schema.sql`
- Run `make run`
