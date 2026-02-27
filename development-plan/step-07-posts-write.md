# Step 7: Posts CRUD — Write Side

## Goal

Add the ability to create, edit, and delete posts. Build form templates, slug generation, database write functions, and the handlers to tie it all together.

## What Was Built

- `handlers.go` — split out from `main.go`; added `handleNewPost`, `handleCreatePost`, `handleEditPost`, `handleUpdatePost`, `handleDeletePost`
- `slug/slug.go` — slug generation from post titles
- `models/queries.go` — `CreatePost()`, `UpdatePost()`, `DeletePost()` functions
- `templates/posts/new.html` — create form
- `templates/posts/edit.html` — edit form
- Updated routes in `main.go`

## File Reorganization

`main.go` was getting long, so handlers were split into `handlers.go`. In Go, multiple files in the same package compile as one unit — `handlers.go` can call `renderTemplate()` defined in `main.go` without any import.

## Slug Generation

```go
// slug/slug.go
package slug

import (
    "regexp"
    "strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)

func Generate(title string) string {
    s := strings.ToLower(title)
    s = strings.ReplaceAll(s, " ", "-")
    s = nonAlphanumeric.ReplaceAllString(s, "")
    s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}
```

`"My First Post!"` → `"my-first-post"`. Replicates what FriendlyID does in the Rails app.

## Write Query Functions

```go
func CreatePost(title, tagline, body, slug string, published bool) (Post, error) {
    var p Post
    err := database.Pool.QueryRow(
        context.Background(),
        `INSERT INTO posts (title, tagline, body, slug, published)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, title, tagline, body, slug, published, created_at, updated_at`,
        title, tagline, body, slug, published,
    ).Scan(
        &p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
        &p.Published, &p.CreatedAt, &p.UpdatedAt,
    )
    if err != nil {
        return Post{}, fmt.Errorf("error creating post: %w", err)
    }
    return p, nil
}
```

`UpdatePost` and `DeletePost` follow the same pattern.

## Form Handling

```go
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    title     := r.FormValue("title")
    tagline   := r.FormValue("tagline")
    body      := r.FormValue("body")
    published := r.FormValue("published") == "true"
    postSlug  := slug.Generate(title)

    post, err := models.CreatePost(title, tagline, body, postSlug, published)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/posts/"+post.Slug, http.StatusSeeOther)
}
```

## Routes

```go
r.Get("/posts/new", handleNewPost)
r.Post("/posts", handleCreatePost)
r.Get("/posts/{slug}/edit", handleEditPost)
r.Post("/posts/{slug}/edit", handleUpdatePost)
r.Post("/posts/{slug}/delete", handleDeletePost)
```

## Concepts Introduced

- **`r.ParseForm()`** — Must be called explicitly before reading form values. Rails does this automatically.
- **`r.FormValue("field")`** — Reads a single form field value. Equivalent to `params[:field]`.
- **Checkbox handling** — HTML checkboxes only send a value when checked, nothing when unchecked. Check for the specific value: `r.FormValue("published") == "true"`.
- **`http.StatusSeeOther` (303)** — The correct redirect status after a POST. Tells the browser to follow up with a GET, preventing form resubmission on refresh. Rails uses this too.
- **`RETURNING` clause** — PostgreSQL returns the created/updated row in a single query. Rails calls `save` then reads the record back in a second query under the hood.
- **`Pool.Exec`** — Used for DELETE where no rows are returned. Contrast with `QueryRow` (single row) and `Query` (multiple rows).
- **Multi-line strings** — Go double-quoted strings cannot span lines. Use backtick raw string literals for multi-line SQL: `` `SELECT ...` ``. Equivalent to Ruby's heredoc (`<<~SQL`).
- **HTML forms and HTTP methods** — HTML only supports `GET` and `POST`. There is no `PUT`/`PATCH`/`DELETE` from a form. Rails works around this with a hidden `_method` field; this app uses `POST /posts/{slug}/edit` and `POST /posts/{slug}/delete`.

## Common Bugs Found

- `RETURNING` clause missing `title` column — caused silent empty fields after create/update
- `DELETE FROM post` (typo) vs `DELETE FROM posts`
- `GET /posts/{slug}/edit` route missing — edit form returned 404
- `handleEditPost` function defined in route but missing from `handlers.go`

## Rails Comparison

| Rails | Go |
|-------|----|
| FriendlyID gem | `slug/slug.go` — hand-rolled slug generator |
| `params[:title]` | `r.FormValue("title")` |
| `@post.save` | `models.CreatePost(...)` |
| `redirect_to @post` | `http.Redirect(w, r, "/posts/"+post.Slug, http.StatusSeeOther)` |
| `render :new` (on error) | `renderTemplate(w, "posts.new", data)` |
| `resources :posts` | Explicit route registration per method |
| `before_action :find_post` | `GetPostBySlug(slug)` called in each handler |
