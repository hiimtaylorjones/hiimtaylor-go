# Step 5: Data Models & Migrations

## Goal

Define the database schema and Go structs for the two core models — `Post` and `Admin` — mirroring the Rails schema but scoped to what this rebuild actually needs.

## What Was Built

- `database/migrations/` — two Goose migration files
- `models/post.go` — `Post` struct
- `models/admin.go` — `Admin` struct

## Key Commands

```bash
brew install goose
goose -dir database/migrations create create_posts sql
goose -dir database/migrations create create_admins sql
goose -dir database/migrations postgres "postgres://localhost:5432/hiimtaylor_go_development?sslmode=disable" up
```

## Why Goose over golang-migrate?

The project started with `golang-migrate`, which splits each migration into two files (`up.sql` and `down.sql`). We switched to **Goose** because it keeps both directions in a single file — closer to Rails' single migration file convention.

```sql
-- +goose Up
CREATE TABLE posts (...);

-- +goose Down
DROP TABLE IF EXISTS posts;
```

Goose tracks applied migrations in a `goose_db_version` table, equivalent to Rails' `schema_migrations`.

## Migration Files

### `create_posts.sql`

```sql
-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    tagline TEXT,
    body TEXT NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS posts;
```

### `create_admins.sql`

```sql
-- +goose Up
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS admins;
```

**Note on `VARCHAR(255)`:** This matches Rails' default — `t.string :column` generates `VARCHAR(255)` in PostgreSQL. Length validation is handled in application code, not the database.

## Model Structs

```go
// models/post.go
type Post struct {
    ID        int
    Title     string
    Tagline   string
    Body      string
    Slug      string
    Published bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

// models/admin.go
type Admin struct {
    ID                int
    Email             string
    EncryptedPassword string
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

## Concepts Introduced

- **Raw SQL migrations** — Unlike Rails' Ruby DSL (`t.string :title`), Go migrations are plain SQL. More verbose, but you see exactly what gets executed.
- **Go structs as models** — No inheritance, no magic. A struct is just a data container. Querying, validation, and associations are functions you write yourself.
- **Uppercase field names** — Go uses capitalization for visibility. Uppercase = exported (accessible from other packages and templates). Lowercase = unexported (private to the package).
- **`REFERENCES ... ON DELETE CASCADE`** — Database-level cascading deletes. Rails often handles this in Ruby with `dependent: :destroy`; here the database does it directly.
- **`SERIAL PRIMARY KEY`** — Auto-incrementing integer ID, same as Rails' default.

## What Was Intentionally Omitted

Compared to the original Rails app, these models were dropped for this rebuild:
- **Comment / Response** — planned for removal from the Rails app too
- **Page** — static pages will be handled via markdown files instead (see Step 10)

## Rails Comparison

| Rails | Go / Goose |
|-------|-----------|
| `rails g migration CreatePosts` | `goose create create_posts sql` |
| `rails db:migrate` | `goose ... up` |
| `schema_migrations` table | `goose_db_version` table |
| `t.string :title` | `title VARCHAR(255) NOT NULL` |
| `class Post < ApplicationRecord` | `type Post struct { ... }` |
| ActiveRecord attributes | Struct fields (must be uppercase) |
| `has_many`, `belongs_to` | Handled in query functions |
