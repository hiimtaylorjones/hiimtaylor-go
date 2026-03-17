# Step 10: Pagination

## Goal

Paginate the posts index to show 10 posts per page, replicating the `will_paginate` behavior from the Rails app — without any gem.

## Status: Complete

## What Was Built

- `models/pagination.go` — `Pagination` struct and `NewPagination()` constructor
- Updated `models/queries.go` — `CountPublishedPosts()` and updated `GetPublishedPosts(page, perPage int)`
- Updated `handlers.go` — `handleListPosts` reads `?page=` param and builds pagination
- Updated `templates/posts/index.html` — prev/next navigation links
- Updated `main.go` — registered `add` and `subtract` template functions via `FuncMap`

## The Pagination Struct

```go
// models/pagination.go
type Pagination struct {
    CurrentPage int
    TotalPages  int
    TotalCount  int
    PerPage     int
    HasPrev     bool
    HasNext     bool
}

func NewPagination(page, perPage, totalCount int) Pagination {
    totalPages := totalCount / perPage
    if totalCount%perPage != 0 {
        totalPages++
    }
    if page < 1 {
        page = 1
    }
    return Pagination{
        CurrentPage: page,
        TotalPages:  totalPages,
        TotalCount:  totalCount,
        PerPage:     perPage,
        HasPrev:     page > 1,
        HasNext:     page < totalPages,
    }
}
```

## Updated Query Functions

```go
func CountPublishedPosts() (int, error) {
    var count int
    err := database.Pool.QueryRow(
        context.Background(),
        `SELECT COUNT(*) FROM posts WHERE published = TRUE`,
    ).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("error counting posts: %w", err)
    }
    return count, nil
}

func GetPublishedPosts(page, perPage int) ([]Post, error) {
    offset := (page - 1) * perPage
    rows, err := database.Pool.Query(
        context.Background(),
        `SELECT id, title, tagline, body, slug, published, created_at, updated_at
         FROM posts WHERE published = TRUE
         ORDER BY created_at DESC
         LIMIT $1 OFFSET $2`,
        perPage, offset,
    )
    // ... scan rows as before
}
```

## Updated Handler

```go
func handleListPosts(w http.ResponseWriter, r *http.Request) {
    const perPage = 10

    page := 1
    if p := r.URL.Query().Get("page"); p != "" {
        if n, err := strconv.Atoi(p); err == nil && n > 0 {
            page = n
        }
    }

    totalCount, err := models.CountPublishedPosts()
    if err != nil {
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }

    posts, err := models.GetPublishedPosts(page, perPage)
    if err != nil {
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }

    pagination := models.NewPagination(page, perPage, totalCount)

    renderTemplate(w, "posts.index", map[string]any{
        "Posts":      posts,
        "Pagination": pagination,
    })
}
```

## Template Functions

Go's template engine doesn't support arithmetic expressions like `{{.CurrentPage - 1}}`. Custom functions are registered via a `FuncMap`:

```go
funcMap := template.FuncMap{
    "add":      func(a, b int) int { return a + b },
    "subtract": func(a, b int) int { return a - b },
}

tmpl := template.New("").Funcs(funcMap)
templates[name] = template.Must(tmpl.ParseFiles(files...))
```

Used in the template as:

```html
{{with .Pagination}}
<nav class="pagination">
    {{if .HasPrev}}
    <a href="/posts?page={{subtract .CurrentPage 1}}">&larr; Newer</a>
    {{end}}
    <span>Page {{.CurrentPage}} of {{.TotalPages}}</span>
    {{if .HasNext}}
    <a href="/posts?page={{add .CurrentPage 1}}">Older &rarr;</a>
    {{end}}
</nav>
{{end}}
```

`{{with .Pagination}}` sets `.` to the `Pagination` struct for the block — cleaner than prefixing every field with `.Pagination`.

## Concepts Introduced

- **`LIMIT` / `OFFSET`** — SQL's pagination primitives. `LIMIT` caps the number of rows returned; `OFFSET` skips the first N rows. `OFFSET = (page - 1) * perPage`.
- **`r.URL.Query().Get("page")`** — Reads a query string parameter. Rails does this automatically with `params[:page]`; in Go you access `r.URL.Query()` explicitly.
- **`strconv.Atoi`** — Converts a string to an integer. Go won't coerce types — you convert explicitly and handle the error.
- **`template.FuncMap`** — A map of custom functions available in templates. The Go equivalent of Rails view helpers. Must be registered before `ParseFiles` is called.
- **`{{with .Pagination}}`** — Scopes `.` to the given value for the block. Equivalent to `<% pagination = @pagination %>` and then using local variables inside.
- **Two queries per page load** — One `COUNT(*)` and one `SELECT` with `LIMIT/OFFSET`. Rails' `will_paginate` does the same thing under the hood.

## Rails Comparison

| Rails | Go |
|-------|----|
| `will_paginate` gem | Hand-rolled `Pagination` struct |
| `Post.paginate(page: params[:page], per_page: 10)` | `GetPublishedPosts(page, perPage)` + `CountPublishedPosts()` |
| `params[:page]` | `r.URL.Query().Get("page")` |
| `will_paginate @posts` helper in view | `{{with .Pagination}}` block in template |
| Rails view helpers (built-in) | `template.FuncMap` (registered manually) |
