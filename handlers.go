package main

import (
	"net/http"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/chi/v5"
	"github.com/hiimtaylorjones/hiimtaylor-go/content"
	"github.com/hiimtaylorjones/hiimtaylor-go/models"
	"github.com/hiimtaylorjones/hiimtaylor-go/slug"
	"github.com/hiimtaylorjones/hiimtaylor-go/uploads"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	html, err := content.Render("home.md")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "home", map[string]any{"Content": html})
}

func handleListPosts(w http.ResponseWriter, r *http.Request) {

	const perPage = 10

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	totalCount, err := models.CountPublishedPosts()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	posts, err := models.GetPublishedPosts(page, perPage)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	pagination := models.NewPagination(page, perPage, totalCount)

	renderTemplate(w, "posts.index", map[string]any{
		"Posts": posts,
		"Pagination": pagination,
	})
}

func handleShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := models.GetPostBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "posts.show", map[string]any{"Post": post})
}

func handleNewPost(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "posts.new", nil)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// Caps form data at 10 MB limit.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	tagline := r.FormValue("tagline")
	body := r.FormValue("body")
	published := r.FormValue("published") == "true"
	postSlug := slug.Generate(title)

	var bannerImageURL string
	file, header, err := r.FormFile("banner_image")
	if err == nil {
		defer file.Close()
		bannerImageURL, err = uploads.Save(file, header)
		if err != nil {
			http.Error(w, "Error saving image", http.StatusInternalServerError)
			return
		}
	}

	post, err := models.CreatePost(title, tagline, body, postSlug, bannerImageURL, published)
	if err != nil {
					http.Error(w, "Error creating post", http.StatusInternalServerError)
					return
	}

	http.Redirect(w, r, "/posts/"+post.Slug, http.StatusSeeOther)
}

func handleEditPost(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "slug")
	post, err := models.GetPostBySlug(postSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "posts.edit", map[string]any{"Post": post})
}

func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "slug")
	post, err := models.GetPostBySlug(postSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
  tagline := r.FormValue("tagline")
  body := r.FormValue("body")
  published := r.FormValue("published") == "true"
	bannerImageURL := post.BannerImageURL

	file, header, err := r.FormFile("banner_image")
	if err == nil {
		defer file.Close()
		bannerImageURL, err = uploads.Save(file, header)
		if err != nil {
			http.Error(w, "Error saving image", http.StatusInternalServerError)
			return
		}
	}

	updated, err := models.UpdatePost(post.ID, title, tagline, body, bannerImageURL, published)
	if err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts/"+updated.Slug, http.StatusSeeOther)
}

func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "slug")
	post, err := models.GetPostBySlug(postSlug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := models.DeletePost(post.ID); err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	html, err := content.Render("resume.md")
	if err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "resume", map[string]any{"Content": html})
}

func handleLoginForm(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
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