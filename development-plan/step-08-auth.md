# Step 8: Admin Authentication

## Goal

Protect write routes behind a login wall, replacing Devise with hand-rolled session-based authentication. The three pillars: password hashing, session management, and auth middleware.

## Status: In Progress

## What Will Be Built

- `middleware/auth.go` — `RequireAdmin` middleware
- `templates/login.html` — login form
- `scripts/seed_admin.go` — one-time admin creation script
- Updated `models/queries.go` — `GetAdminByEmail()`, `CreateAdmin()`
- Updated `handlers.go` — `handleLoginForm`, `handleLogin`, `handleLogout`
- Updated `main.go` — session manager, protected route group

## Key Commands

```bash
go get golang.org/x/crypto/bcrypt
go get github.com/alexedwards/scs/v2
go mod vendor
```

## Session Manager Setup

```go
// main.go
var sessionManager *scs.SessionManager

func main() {
    sessionManager = scs.New()
    sessionManager.Lifetime = 24 * time.Hour

    // Wrap the entire router with session middleware
    log.Fatal(http.ListenAndServe(":3000", sessionManager.LoadAndSave(r)))
}
```

## Auth Middleware

```go
// middleware/auth.go
func RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        adminID := sessionManager.GetString(r, "admin_id")
        if adminID == "" {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Protected Route Group

```go
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
```

## Seeding an Admin

```go
// scripts/seed_admin.go
// go:build ignore

func main() {
    database.Connect()
    hash, _ := bcrypt.GenerateFromPassword([]byte("your-password"), bcrypt.DefaultCost)
    admin, _ := models.CreateAdmin("admin@hiimtaylorjones.com", string(hash))
    fmt.Printf("Created admin: %s\n", admin.Email)
}
```

Run with: `go run scripts/seed_admin.go`

The `//go:build ignore` tag prevents this file from being included in regular builds.

## Concepts Introduced

- **`bcrypt.GenerateFromPassword`** — Hashes a plaintext password with a cost factor. Devise uses bcrypt internally too.
- **`bcrypt.CompareHashAndPassword`** — Verifies a plaintext password against a stored hash without ever decrypting it.
- **Same error for wrong email and wrong password** — Prevents user enumeration attacks. Never reveal which part was incorrect.
- **`sessionManager.Put(r.Context(), key, value)`** — Stores data in the session. Rails does this with `sign_in(admin)` which ultimately writes to the session.
- **`sessionManager.LoadAndSave(r)`** — Wraps the router so every request loads its session from the store and saves it back after the handler runs.
- **`r.Group()`** — Chi sub-router that inherits parent middleware and adds its own. The Go equivalent of scoping `before_action` to specific controller actions.
- **`//go:build ignore`** — Build tag that excludes a file from regular compilation. Used for one-off scripts that shouldn't be part of the binary.

## Rails Comparison

| Rails (Devise) | Go |
|----------------|----|
| `devise :database_authenticatable` | `bcrypt.CompareHashAndPassword` |
| `before_action :authenticate_admin!` | `r.Use(authmiddleware.RequireAdmin)` |
| `sign_in(admin)` | `sessionManager.Put(r.Context(), "admin_id", ...)` |
| `sign_out(admin)` | `sessionManager.Remove(r.Context(), "admin_id")` |
| `current_admin` | Look up admin by ID stored in session |
| `admin_signed_in?` | `sessionManager.GetString(r, "admin_id") != ""` |
| `db/seeds.rb` | `scripts/seed_admin.go` |
| Session stored in cookie | Session ID in cookie, data server-side (scs default) |
