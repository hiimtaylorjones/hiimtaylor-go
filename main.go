  package main

  import (
        "html/template"
        "log"
        "net/http"
        "path/filepath"

        "github.com/go-chi/chi/v5"
        "github.com/go-chi/chi/v5/middleware"
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
      r.Get("/posts", handleListPosts)
      r.Get("/posts/{slug}", handleShowPost)
      r.Get("/resume", handleResume)

      log.Println("Server starting on http://localhost:3000")
      log.Fatal(http.ListenAndServe(":3000", r))
  }

  func handleHome(w http.ResponseWriter, r *http.Request) {
      renderTemplate(w, "home", nil)
  }

  func handleListPosts(w http.ResponseWriter, r *http.Request) {
      // Short term solution until we implement database connection
      data := map[string]any{
            "Posts": []map[string]string{
                  {"Title": "My First Post", "Tagline": "Hello world", "Slug": "my-first-post"},
                  {"Title": "Learning Go", "Tagline": "Coming from Ruby", "Slug": "learning-go"},
            },
      }
      renderTemplate(w, "posts.index", data)
  }

  func handleShowPost(w http.ResponseWriter, r *http.Request) {
      slug := chi.URLParam(r, "slug")
      
      // Another temporary solution until database implementation.
      data := map[string]any{
            "Post": map[string]string{
                  "Title": slug,
                  "Tagline": "A placeholder post",
                  "Body": "This is the body of the post.",
            },
      }
      renderTemplate(w, "posts.show", data)
  }

  func handleResume(w http.ResponseWriter, r *http.Request) {
      renderTemplate(w, "resume", nil)
  }