# Step 2: Add Routing with Chi

## Goal

Replace Go's default mux with a real router that supports URL parameters, HTTP method-specific routing, and middleware.

## What Was Built

- Replaced `http.DefaultServeMux` with `chi.NewRouter()`
- Four named routes mirroring the core pages of the Rails site
- Request logging and panic recovery middleware

## Key Commands

```bash
go get github.com/go-chi/chi/v5
```

## The Code

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Get("/", handleHome)
    r.Get("/posts", handleListPosts)
    r.Get("/posts/{slug}", handleShowPost)
    r.Get("/resume", handleResume)

    log.Fatal(http.ListenAndServe(":3000", r))
}

func handleShowPost(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    fmt.Fprintf(w, "Showing post: %s", slug)
}
```

## Concepts Introduced

- **`chi.NewRouter()`** — Creates a router instance with method-specific routing and URL parameter support.
- **`r.Use(...)`** — Registers middleware that runs on every request. Middleware wraps handlers in a chain.
- **`middleware.Logger`** — Logs each request with method, path, status, and duration. Equivalent to Rails' development request log.
- **`middleware.Recoverer`** — Catches panics and returns a 500 instead of crashing the server.
- **`r.Get("/posts/{slug}", ...)`** — Method-specific route with a URL parameter. Chi uses `{param}` syntax; Rails uses `:param`.
- **`chi.URLParam(r, "slug")`** — Extracts the URL parameter value. Equivalent to `params[:slug]` in Rails.
- **Named handler functions** — Each route gets its own function with the signature `func(http.ResponseWriter, *http.Request)`.

## Rails Comparison

| Rails | Go / Chi |
|-------|----------|
| `config/routes.rb` | `chi.NewRouter()` in `main.go` |
| `get '/posts/:slug'` | `r.Get("/posts/{slug}", ...)` |
| `params[:slug]` | `chi.URLParam(r, "slug")` |
| `before_action` (global) | `r.Use(middleware)` |
| Rails request log | `middleware.Logger` |
| Rails exception handling | `middleware.Recoverer` |
