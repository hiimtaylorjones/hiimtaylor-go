package main

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/hiimtaylorjones/hiimtaylor-go/database"
)

var templates map[string]*template.Template

func loadTemplates() {
    templates = make(map[string]*template.Template)

    layouts, _ := filepath.Glob("templates/layouts/*.html")
    partials, _ := filepath.Glob("templates/partials/*.html")

    pages := map[string]string{
        "home":           "templates/home.html",
        "resume":         "templates/resume.html",
        "posts.index":    "templates/posts/index.html",
        "posts.show":     "templates/posts/show.html",
        "posts.new":      "templates/posts/new.html",
        "posts.edit":     "templates/posts/edit.html",
    }

    for name, pages := range pages {
        files := append(layouts, partials...)
        files = append(files, pages)
        templates[name] = template.Must(template.ParseFiles(files...))
    }
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
    tmpl, ok := templates[name]
    if !ok {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }
    err := tmpl.ExecuteTemplate(w, "base", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    database.Connect()
    defer database.Close()

    loadTemplates()

    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Fetch and Serve Static Files
    fileServer := http.FileServer(http.Dir("static"))
    r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

    // Routes
    r.Get("/", handleHome)
    r.Get("/posts/new", handleNewPost)
    r.Get("/posts", handleListPosts)
    r.Post("/posts", handleCreatePost)
    r.Get("/posts/{slug}", handleShowPost)
    r.Get("/posts/{slug}/edit", handleEditPost)
    r.Post("/posts/{slug}/edit", handleUpdatePost)
    r.Post("/posts/{slug}/delete", handleDeletePost)
    r.Get("/resume", handleResume)

    log.Println("Server starting on http://localhost:3000")
    log.Fatal(http.ListenAndServe(":3000", r))
}
