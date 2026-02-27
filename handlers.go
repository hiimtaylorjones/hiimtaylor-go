package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hiimtaylorjones/hiimtaylor-go/models"
	"github.com/hiimtaylorjones/hiimtaylor-go/slug"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", nil)
}

func handleListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetPublishedPosts()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "posts.index", map[string]any{"Posts": posts})
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	tagline := r.FormValue("tagline")
	body := r.FormValue("body")
	published := r.FormValue("published") == "true"
	postSlug := slug.Generate(title)

	post, err := models.CreatePost(title, tagline, body, postSlug, published)
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
  tagline := r.FormValue("tagline")
  body := r.FormValue("body")
  published := r.FormValue("published") == "true"

	updated, err := models.UpdatePost(post.ID, title, tagline, body, published)
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
	renderTemplate(w, "resume", nil)
}