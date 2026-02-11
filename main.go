  package main

  import (
        "fmt"
        "log"
        "net/http"

        "github.com/go-chi/chi/v5"
        "github.com/go-chi/chi/v5/middleware"
  )

  func main() {
        r := chi.NewRouter()

        // Middleware
        r.Use(middleware.Logger)
        r.Use(middleware.Recoverer)

        // Routes
        r.Get("/", handleHome)
        r.Get("/posts", handleListPosts)
        r.Get("/posts/{slug}", handleShowPost)
        r.Get("/resume", handleResume)

        log.Println("Server starting on http://localhost:3000")
        log.Fatal(http.ListenAndServe(":3000", r))
  }

  func handleHome(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to hiimtaylorjones.com")
  }

  func handleListPosts(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "All posts will go here")
  }

  func handleShowPost(w http.ResponseWriter, r *http.Request) {
        slug := chi.URLParam(r, "slug")
        fmt.Fprintf(w, "Showing post: %s", slug)
  }

  func handleResume(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Resume will go here")
  }