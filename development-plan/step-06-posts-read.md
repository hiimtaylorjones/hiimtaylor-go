# Step 6: Posts CRUD — Read Side

## Goal

Replace placeholder data in handlers with real database queries. Render post bodies as HTML from markdown. Wire everything through proper structs instead of ad-hoc maps.

## What Was Built

- `models/queries.go` — `GetPublishedPosts()` and `GetPostBySlug()` query functions
- `models/post.go` — `RenderedBody()` method using Goldmark for markdown-to-HTML
- Updated `handlers.go` — `handleListPosts` and `handleShowPost` use real data
- Updated `templates/posts/show.html` — renders `{{.Post.RenderedBody}}` instead of raw body text

## Key Commands

```bash
go get github.com/yuin/goldmark
```

## Query Functions

```go
func GetPublishedPosts() ([]Post, error) {
    rows, err := database.Pool.Query(
        context.Background(),
        `SELECT id, title, tagline, body, slug, published, created_at, updated_at
         FROM posts WHERE published = TRUE ORDER BY created_at DESC`,
    )
    if err != nil {
        return nil, fmt.Errorf("error querying posts: %w", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var p Post
        err := rows.Scan(
            &p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
            &p.Published, &p.CreatedAt, &p.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning post: %w", err)
        }
        posts = append(posts, p)
    }
    return posts, nil
}

func GetPostBySlug(slug string) (Post, error) {
    var p Post
    err := database.Pool.QueryRow(
        context.Background(),
        `SELECT id, title, tagline, body, slug, published, created_at, updated_at
         FROM posts WHERE slug = $1`,
        slug,
    ).Scan(
        &p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
        &p.Published, &p.CreatedAt, &p.UpdatedAt,
    )
    if err != nil {
        return Post{}, fmt.Errorf("post not found: %w", err)
    }
    return p, nil
}
```

## Markdown Rendering

```go
func (p Post) RenderedBody() template.HTML {
    var buf bytes.Buffer
    if err := goldmark.Convert([]byte(p.Body), &buf); err != nil {
        return template.HTML(p.Body)
    }
    return template.HTML(buf.String())
}
```

Used in the template as `{{.Post.RenderedBody}}`.

## Updated Handlers

```go
func handleListPosts(w http.ResponseWriter, r *http.Request) {
    posts, err := models.GetPublishedPosts()
    if err != nil {
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "posts.index", map[string]any{"Posts": posts})
}

func handleShowPost(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    post, err := models.GetPostBySlug(slug)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    renderTemplate(w, "posts.show", map[string]any{"Post": post})
}
```

## Concepts Introduced

- **`rows.Scan(...)`** — Manually maps each column to a struct field. The order must exactly match the SELECT clause. ActiveRecord does this automatically by matching column names to attribute names.
- **`$1`, `$2`** — PostgreSQL parameterized query placeholders (Rails uses `?` or named params). Prevents SQL injection.
- **`defer rows.Close()`** — Ensures the result set is released back to the pool after the function returns. No equivalent in Rails — ActiveRecord handles this transparently.
- **`fmt.Errorf("...: %w", err)`** — Wraps errors with context. The `%w` verb allows callers to inspect the original error with `errors.Is()`. Go handles errors as return values, not exceptions.
- **`template.HTML`** — Marks a string as safe HTML, bypassing the template engine's auto-escaping. Without this, markdown-rendered `<p>` tags would be printed as literal text. Equivalent to `.html_safe` or `raw()` in Rails.
- **Method on a struct** — `func (p Post) RenderedBody()` is Go's equivalent of an instance method. Called as `post.RenderedBody()` in templates.
- **`Pool.Query` vs `Pool.QueryRow` vs `Pool.Exec`** — Three variants: `Query` for multiple rows, `QueryRow` for a single row, `Exec` for no rows returned (INSERT/UPDATE/DELETE without RETURNING).

## Common Gotchas Encountered

- **Mixed indentation** — Go is strict about tabs vs spaces. `gofmt` enforces tabs.
- **`p.body` vs `p.Body`** — Lowercase struct fields are unexported and inaccessible outside the package. Always use uppercase for fields that templates or other packages need.
- **Trailing comma rule** — When a function call spans multiple lines, the last argument needs a trailing comma before the closing `)`.

## Rails Comparison

| Rails | Go |
|-------|----|
| `Post.where(published: true).order(created_at: :desc)` | `GetPublishedPosts()` with raw SQL |
| `Post.find_by!(slug: slug)` | `GetPostBySlug(slug)` |
| `@post.body.html_safe` | `post.RenderedBody()` returning `template.HTML` |
| Commonmarker gem | `github.com/yuin/goldmark` |
| `rescue ActiveRecord::RecordNotFound` | `if err != nil { http.NotFound(...) }` |
| ActiveRecord lazy loading | Eager — data is fetched when the function is called |
