# Step 11: Image Uploads

## Goal

Support banner images on posts, replicating the Active Storage + S3 setup from the Rails app — starting with local disk storage in `static/uploads/`.

## Status: Complete

## What Was Built

- `database/migrations/XXXXXX_add_banner_image_to_posts.sql` — adds `banner_image_url` column
- `uploads/uploads.go` — saves multipart file uploads to disk, returns public URL
- Updated `models/post.go` — added `BannerImageURL` field to `Post` struct
- Updated `models/queries.go` — all post queries updated to include `banner_image_url`
- Updated `handlers.go` — `handleCreatePost` and `handleUpdatePost` use `ParseMultipartForm` and `uploads.Save()`
- Updated `templates/posts/new.html` and `templates/posts/edit.html` — file input with `enctype="multipart/form-data"`
- Updated `templates/posts/show.html` — conditionally renders banner image

## Migration

```sql
-- +goose Up
ALTER TABLE posts ADD COLUMN banner_image_url VARCHAR(500);

-- +goose Down
ALTER TABLE posts DROP COLUMN banner_image_url;
```

## The Uploads Package

```go
package uploads

const uploadDir = "static/uploads"

func Save(file multipart.File, header *multipart.FileHeader) (string, error) {
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        return "", fmt.Errorf("could not create upload directory: %w", err)
    }

    ext := filepath.Ext(header.Filename)
    filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    dest := filepath.Join(uploadDir, filename)

    out, err := os.Create(dest)
    if err != nil {
        return "", fmt.Errorf("could not create file: %w", err)
    }
    defer out.Close()

    if _, err := io.Copy(out, file); err != nil {
        return "", fmt.Errorf("could not save file: %w", err)
    }

    return "/static/uploads/" + filename, nil
}
```

## Form Handling

```go
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    // ... read text fields ...

    var bannerImageURL string
    file, header, err := r.FormFile("banner_image")
    if err == nil { // file was uploaded
        defer file.Close()
        bannerImageURL, err = uploads.Save(file, header)
        if err != nil {
            http.Error(w, "Error saving image", http.StatusInternalServerError)
            return
        }
    }
    // if err != nil, no file was uploaded — bannerImageURL stays ""

    post, err := models.CreatePost(title, tagline, body, postSlug, bannerImageURL, published)
    // ...
}
```

For `handleUpdatePost`, the existing `BannerImageURL` is used as the default — only replaced if a new file is uploaded:

```go
bannerImageURL := post.BannerImageURL // keep existing by default

file, header, err := r.FormFile("banner_image")
if err == nil {
    defer file.Close()
    bannerImageURL, err = uploads.Save(file, header)
    // ...
}
```

## Template Changes

Form tags need `enctype="multipart/form-data"` — without it, file data is never sent:

```html
<form method="POST" action="/posts" enctype="multipart/form-data">
```

Post show template conditionally renders the banner:

```html
{{if .Post.BannerImageURL}}
<img src="{{.Post.BannerImageURL}}" alt="{{.Post.Title}}" class="banner-image">
{{end}}
```

## Concepts Introduced

- **`r.ParseMultipartForm(10 << 20)`** — Parses a multipart form with a max memory limit of 10MB (`10 << 20` is bitwise left shift, a common Go idiom for byte sizes). `r.ParseForm()` does not handle file uploads.
- **`r.FormFile("field")`** — Returns the uploaded file, its metadata header, and an error. If no file was submitted, `err` is non-nil — we use this to make the upload optional rather than required.
- **`io.Copy(out, file)`** — Streams the uploaded file to disk without loading it entirely into memory. Rails' Active Storage does this same streaming internally.
- **`time.Now().UnixNano()`** — Generates a unique filename from the current nanosecond timestamp. Active Storage uses a random token for the same purpose.
- **`filepath.Ext`** — Extracts the file extension (`.jpg`, `.png`) from the original filename to preserve it on the saved file.
- **`os.MkdirAll`** — Creates the upload directory and any missing parents. Safe to call even if the directory already exists (unlike `os.Mkdir`).
- **`enctype="multipart/form-data"`** — Required HTML attribute for any form that submits files. Without it, the browser sends files as plain text and the upload is lost.

## Common Bugs Found

- `GetPublishedPosts` Scan had `Published` and `BannerImageURL` swapped relative to the SELECT column order — columns must be scanned in exactly the same order they appear in the SELECT clause
- `CreatePost` was missing `bannerImageURL` from the args list and `&p.BannerImageURL` from Scan
- `UpdatePost` was missing `banner_image_url=$5` in the SET clause, `id=$5` should have been `id=$6`, and `banner_image_url` was absent from RETURNING and Scan

## Rails Comparison

| Rails | Go |
|-------|----|
| Active Storage + `has_one_attached :banner_image` | `uploads.Save()` writing to `static/uploads/` |
| Direct S3 upload via Active Storage | Local disk write via `io.Copy` (S3 is a future swap) |
| `banner_image.attached?` | `{{if .Post.BannerImageURL}}` |
| `url_for(@post.banner_image)` | `{{.Post.BannerImageURL}}` (stored as a plain path) |
| `params[:post][:banner_image]` | `r.FormFile("banner_image")` |
| Active Storage migration auto-generated | Manual `ALTER TABLE` migration |
| Stores file metadata in `active_storage_blobs` | Stores URL string directly in `posts.banner_image_url` |
