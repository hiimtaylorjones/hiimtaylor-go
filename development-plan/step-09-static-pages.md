# Step 9: Static Pages (Home & Resume)

## Goal

Serve the home and resume pages from markdown files on disk, replicating the Rails app's `.html.md` view convention without a database. Content lives in files, not rows.

## Status: Complete

## What Was Built

- `content/content.go` — package for reading and rendering markdown files to HTML
- `content/home.md` — home page content
- `content/resume.md` — resume content
- Updated `templates/home.html` and `templates/resume.html` — render `{{.Content}}` instead of hardcoded text
- Updated `handleHome` and `handleResume` in `handlers.go` to use `content.Render()`

## Directory Structure

```
content/
  content.go    <- Render() helper
  home.md       <- home page copy
  resume.md     <- resume/CV content
```

## The Content Package

```go
package content

import (
    "bytes"
    "fmt"
    "html/template"
    "os"

    "github.com/yuin/goldmark"
)

func Render(filename string) (template.HTML, error) {
    raw, err := os.ReadFile("content/" + filename)
    if err != nil {
        return "", fmt.Errorf("could not read content file: %w", err)
    }

    var buf bytes.Buffer
    if err := goldmark.Convert(raw, &buf); err != nil {
        return "", fmt.Errorf("could not render markdown: %w", err)
    }

    return template.HTML(buf.String()), nil
}
```

## Updated Templates

```html
<!-- templates/home.html -->
{{define "content"}}
<div class="page-content">
    {{.Content}}
</div>
{{end}}
```

## Updated Handlers

```go
func handleHome(w http.ResponseWriter, r *http.Request) {
    html, err := content.Render("home.md")
    if err != nil {
        http.Error(w, "Could not load page", http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "home", map[string]any{"Content": html})
}

func handleResume(w http.ResponseWriter, r *http.Request) {
    html, err := content.Render("resume.md")
    if err != nil {
        http.Error(w, "Could not load page", http.StatusInternalServerError)
        return
    }
    renderTemplate(w, "resume", map[string]any{"Content": html})
}
```

## Concepts Introduced

- **`os.ReadFile`** — Reads a file from disk into a byte slice. Simple, no streaming. The Go equivalent of `File.read` in Ruby.
- **File-based content vs. database content** — Not everything belongs in a database. Static pages you edit infrequently are a good fit for version-controlled markdown files. Rails supports this too with `.html.md` views, but it's more magical.
- **Live edits without restart** — Because files are read on each request (no caching yet), editing a markdown file takes effect on the next page load without restarting the server. Templates, by contrast, are parsed once at startup.
- **`template.HTML` return type** — Marks the rendered HTML as safe, preventing the template engine from double-escaping it. Same pattern used in `RenderedBody()` on the Post model in Step 6.
- **Separation of content and structure** — The template defines layout and CSS classes; the markdown file holds the actual prose. A clean separation that makes content easy to update without touching Go code.

## Rails Comparison

| Rails | Go |
|-------|----|
| `app/views/basic_page/resume.html.md` | `content/resume.md` |
| Custom `MarkdownHandler` initializer | `content.Render()` function |
| Rails autoloads and renders `.md` views | Handler explicitly reads file and converts |
| `render :resume` | `content.Render("resume.md")` + `renderTemplate(...)` |
| `html_safe` | `template.HTML` return type |
