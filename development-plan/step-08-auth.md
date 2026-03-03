# Step 8: Admin Authentication

## Goal

Protect write routes behind a login wall, replacing Devise with hand-rolled session-based authentication. The three pillars: password hashing, session management, and auth middleware.

## Status: Complete

## What Was Built

- `middleware/auth.go` — `RequireAdmin` middleware and session manager wiring
- `templates/login.html` — login form with error display
- `scripts/seed_admin.go` — one-time admin creation script
- Updated `models/queries.go` — `GetAdminByEmail()` and `CreateAdmin()` query functions
- Updated `handlers.go` — `handleLoginForm`, `handleLogin`, `handleLogout`
- Updated `main.go` — session manager setup, protected route group, `sessionManager.LoadAndSave` wrapping the router

## Key Commands

```bash
go get golang.org/x/crypto/bcrypt
go get github.com/alexedwards/scs/v2
go mod vendor

# Run once to create the admin account
go run scripts/seed_admin.go
```

## Auth Middleware

```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func SetSessionManager(sm *scs.SessionManager) {
    sessionManager = sm
}

func RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        adminID := sessionManager.GetString(r.Context(), "admin_id")
        if adminID == "" {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Session Manager Setup in `main.go`

```go
var sessionManager *scs.SessionManager

func main() {
    // ...
    sessionManager = scs.New()
    sessionManager.Lifetime = 6 * time.Hour
    authmiddleware.SetSessionManager(sessionManager)

    // ...

    // Wrap router with session middleware
    log.Fatal(http.ListenAndServe(":3000", sessionManager.LoadAndSave(r)))
}
```

## Protected Route Group

```go
// Public routes
r.Get("/", handleHome)
r.Get("/posts", handleListPosts)
r.Get("/posts/{slug}", handleShowPost)
r.Get("/resume", handleResume)

// Auth routes
r.Get("/login", handleLoginForm)
r.Post("/login", handleLogin)
r.Post("/logout", handleLogout)

// Protected routes — admin only
r.Group(func(r chi.Router) {
    r.Use(authmiddleware.RequireAdmin)
    r.Get("/posts/new", handleNewPost)
    r.Post("/posts", handleCreatePost)
    r.Get("/posts/{slug}/edit", handleEditPost)
    r.Post("/posts/{slug}/edit", handleUpdatePost)
    r.Post("/posts/{slug}/delete", handleDeletePost)
})
```

## Login Handler

```go
func handleLogin(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email    := r.FormValue("email")
    password := r.FormValue("password")

    admin, err := models.GetAdminByEmail(email)
    if err != nil {
        renderTemplate(w, "login", map[string]any{"Error": "Invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(admin.EncryptedPassword), []byte(password)); err != nil {
        renderTemplate(w, "login", map[string]any{"Error": "Invalid email or password"})
        return
    }

    sessionManager.Put(r.Context(), "admin_id", fmt.Sprintf("%d", admin.ID))
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
    sessionManager.Remove(r.Context(), "admin_id")
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
```

## Admin Query Functions

```go
func GetAdminByEmail(email string) (Admin, error) {
    var a Admin
    err := database.Pool.QueryRow(
        context.Background(),
        `SELECT id, email, encrypted_password FROM admins WHERE email = $1`,
        email,
    ).Scan(&a.ID, &a.Email, &a.EncryptedPassword)
    if err != nil {
        return Admin{}, fmt.Errorf("admin not found: %w", err)
    }
    return a, nil
}
```

## Seeding an Admin

```go
//go:build ignore

package main

func main() {
    database.Connect()
    defer database.Close()

    password := "your-password-here"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    admin, _ := models.CreateAdmin("admin@hiimtaylorjones.com", string(hash))
    fmt.Printf("Created admin: %s (id: %d)\n", admin.Email, admin.ID)
}
```

Run with: `go run scripts/seed_admin.go`

The `//go:build ignore` tag prevents this file from being included in regular `go build` runs.

## Concepts Introduced

- **`bcrypt.GenerateFromPassword`** — Hashes a plaintext password with a cost factor. Devise uses bcrypt under the hood too — this is the same algorithm.
- **`bcrypt.CompareHashAndPassword`** — Verifies a plaintext password against a stored hash without ever storing or decrypting the plaintext.
- **Same error for wrong email and wrong password** — Prevents user enumeration attacks. Never reveal which part of the credentials was incorrect.
- **`scs.New()`** — Creates a session manager. By default, sessions are stored server-side with the session ID in a cookie. The `Lifetime` controls how long a session lasts before expiring.
- **`sessionManager.LoadAndSave(r)`** — Wraps the entire router so every request loads its session data and saves any changes after the handler runs.
- **`sessionManager.Put(r.Context(), key, value)`** — Stores a value in the session. Rails does this with `sign_in(admin)` which ultimately writes to the session too.
- **`sessionManager.Remove(r.Context(), key)`** — Removes a value from the session, effectively logging the user out.
- **`r.Group()`** — Chi sub-router that inherits parent middleware and adds its own. The Go equivalent of scoping `before_action` to specific controller actions.
- **`authmiddleware` import alias** — We alias our middleware package to `authmiddleware` to avoid a name clash with Chi's built-in `middleware` package.
- **`//go:build ignore`** — Build tag that excludes a file from regular compilation. Used for one-off scripts that shouldn't be compiled into the main binary.

## Rails Comparison

| Rails (Devise) | Go |
|----------------|----|
| `devise :database_authenticatable` | `bcrypt.CompareHashAndPassword` |
| `before_action :authenticate_admin!` | `r.Use(authmiddleware.RequireAdmin)` |
| `sign_in(admin)` | `sessionManager.Put(r.Context(), "admin_id", ...)` |
| `sign_out(admin)` | `sessionManager.Remove(r.Context(), "admin_id")` |
| `current_admin` | Look up admin by ID stored in session |
| `admin_signed_in?` | `sessionManager.GetString(r.Context(), "admin_id") != ""` |
| `db/seeds.rb` | `scripts/seed_admin.go` with `//go:build ignore` |
| Session cookie + server-side store | scs cookie + server-side store (same model) |
| Devise generates all login views and routes | Manually wired login form, handlers, and routes |
