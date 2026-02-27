# Step 1: Initialize the Go Module & Hello World Server

## Goal

Get a running HTTP server that responds on `localhost:3000` — the Go equivalent of `rails new` and `rails server`.

## What Was Built

- `go.mod` — module definition file
- `main.go` — a minimal HTTP server using only the standard library

## Key Commands

```bash
go mod init github.com/hiimtaylorjones/hiimtaylor-go
go run .
```

## The Code

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    log.Println("Server starting on http://localhost:3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}
```

## Concepts Introduced

- **`package main`** — Every Go executable needs a `main` package. This is the entry point.
- **`import`** — Go has no autoloading. Every dependency is declared explicitly.
- **`http.HandleFunc`** — Registers a route handler. The handler receives a `ResponseWriter` (where you write the response) and a `*Request` (the incoming request).
- **`fmt.Fprintf(w, ...)`** — Writes text to the response. Equivalent to `render plain:` in Rails.
- **`http.ListenAndServe(":3000", nil)`** — Starts the server. The `nil` uses Go's default router (`DefaultServeMux`).
- **`log.Fatal`** — Logs and exits immediately if the server fails to start.
- **`go run .`** — Compiles and runs in one step. `go build .` produces a standalone binary.

## Rails Comparison

| Rails | Go |
|-------|----|
| `rails new` | `go mod init` |
| `rails server` | `go run .` |
| Puma + Rack + Rails framework | `net/http` (stdlib) |
| `render plain: "Hello"` | `fmt.Fprintf(w, "Hello")` |
| `bin/rails` entry point | `func main()` |
