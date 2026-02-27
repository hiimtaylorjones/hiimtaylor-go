# Step 4: Connect to PostgreSQL

## Goal

Connect the Go app to a PostgreSQL database using a connection pool, with environment-based configuration — the Go equivalent of `database.yml` and ActiveRecord's connection management.

## What Was Built

- `database/database.go` — connection pool setup, connect/close functions
- Wired `database.Connect()` and `defer database.Close()` into `main()`

## Key Commands

```bash
createdb hiimtaylor_go_development
go get github.com/jackc/pgx/v5/pgxpool
```

## The Code

```go
package database

import (
    "context"
    "log"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "postgres://localhost:5432/hiimtaylor_go_development?sslmode=disable"
    }

    var err error
    Pool, err = pgxpool.New(context.Background(), connStr)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }

    err = Pool.Ping(context.Background())
    if err != nil {
        log.Fatalf("Unable to ping database: %v\n", err)
    }

    log.Println("Connected to database")
}

func Close() {
    Pool.Close()
}
```

Wired into `main()`:

```go
func main() {
    database.Connect()
    defer database.Close()
    // ...
}
```

## Concepts Introduced

- **`pgxpool`** — A connection pool for PostgreSQL. Rails manages this automatically via ActiveRecord's `pool` setting in `database.yml`. In Go, you create and manage it explicitly.
- **`os.Getenv("DATABASE_URL")`** — Reads config from the environment. Falls back to a local default for development. Equivalent to how Rails reads `DB_USERNAME`/`DB_PASSWORD` env vars in production.
- **`context.Background()`** — Go uses `context` to manage timeouts and cancellation for I/O operations. `context.Background()` means "no deadline, no cancellation." You'll see it on every database call.
- **`Pool.Ping()`** — A health-check query to verify the connection is alive on startup.
- **`log.Fatalf`** — Logs a formatted error and exits. If the database is unreachable, there's no point starting the server.
- **`defer database.Close()`** — `defer` schedules a function to run when the surrounding function returns. Ensures the connection pool is cleaned up on shutdown. Similar to Ruby's `at_exit`.
- **Exported `Pool` variable** — Uppercase `P` makes it accessible from other packages (`models`, etc.). The Go equivalent of ActiveRecord's global connection.

## Notes on `log.Fatal` vs `log.Fatalf`

A common gotcha: `log.Fatal` does **not** format strings — it prints them literally. Use `log.Fatalf` (with `f`) when you want `%v`, `%s`, or other format verbs.

```go
log.Fatal("error: %v", err)   // Wrong — prints literally "%v"
log.Fatalf("error: %v", err)  // Correct — interpolates the error
```

## Rails Comparison

| Rails | Go |
|-------|----|
| `database.yml` | `os.Getenv("DATABASE_URL")` with fallback |
| ActiveRecord connection pool | `pgxpool.New()` |
| `pool: 5` in database.yml | pgxpool manages this automatically |
| `ActiveRecord::Base.connection` | `database.Pool` |
| `rails db:create` | `createdb <dbname>` (psql CLI) |
