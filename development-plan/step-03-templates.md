# Step 3: HTML Templates & Layouts

## Goal

Render real HTML pages using Go's `html/template` package, replicating the Rails layout system of a base layout, shared partials, and per-page content blocks.

## What Was Built

- `templates/layouts/base.html` — master layout (equivalent to `application.html.erb`)
- `templates/partials/header.html` — navigation partial
- `templates/partials/footer.html` — footer partial
- `templates/home.html`, `templates/resume.html` — static page templates
- `templates/posts/index.html`, `templates/posts/show.html` — post templates
- `static/css/style.css` — base stylesheet
- `loadTemplates()` — parses and caches all templates at startup
- `renderTemplate()` — shared helper for executing templates

## Directory Structure

```
templates/
  layouts/
    base.html
  partials/
    header.html
    footer.html
  home.html
  resume.html
  posts/
    index.html
    show.html
static/
  css/
    style.css
```

## Key Concepts

### Template Composition

Go templates use `{{define}}` and `{{template}}` to compose layouts:

```html
<!-- layouts/base.html -->
{{define "base"}}
<html>
  <body>
    {{template "header" .}}
    <main>{{template "content" .}}</main>
    {{template "footer" .}}
  </body>
</html>
{{end}}

<!-- posts/index.html -->
{{define "content"}}
<h1>Posts</h1>
{{range .Posts}}
  <h2>{{.Title}}</h2>
{{end}}
{{end}}
```

### Template Parsing

Templates must be explicitly parsed and combined at startup:

```go
var templates map[string]*template.Template

func loadTemplates() {
    templates = make(map[string]*template.Template)
    layouts, _ := filepath.Glob("templates/layouts/*.html")
    partials, _ := filepath.Glob("templates/partials/*.html")

    pages := map[string]string{
        "posts.index": "templates/posts/index.html",
        // ...
    }

    for name, page := range pages {
        files := append(layouts, partials...)
        files = append(files, page)
        templates[name] = template.Must(template.ParseFiles(files...))
    }
}
```

### Rendering

```go
func renderTemplate(w http.ResponseWriter, name string, data any) {
    tmpl := templates[name]
    tmpl.ExecuteTemplate(w, "base", data)
}
```

### Serving Static Files

```go
fileServer := http.FileServer(http.Dir("static"))
r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
```

## Template Syntax Reference

| Rails (ERB) | Go template |
|------------|-------------|
| `<%= @post.title %>` | `{{.Post.Title}}` |
| `<% @posts.each do \|p\| %>` | `{{range .Posts}}` |
| `<% end %>` | `{{end}}` |
| `<%= yield %>` | `{{template "content" .}}` |
| `<%= render "partials/header" %>` | `{{template "header" .}}` |
| `<% if admin_signed_in? %>` | `{{if .IsAdmin}}` |

## Rails Comparison

| Rails | Go |
|-------|----|
| `application.html.erb` | `layouts/base.html` with `{{define "base"}}` |
| `<%= yield %>` | `{{template "content" .}}` |
| `<%= render partial %>` | `{{template "partial_name" .}}` |
| `@posts` instance variable | Data passed as third arg to `renderTemplate` |
| View lookup by convention | Explicit map in `loadTemplates()` |
| `html_safe` / `raw()` | `template.HTML` type |
| `<%= @posts.each %>` | `{{range .Posts}}` |
