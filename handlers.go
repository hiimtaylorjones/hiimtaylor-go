package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hiimtaylorjones/hiimtaylor-go/models"
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

func handleResume(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "resume", nil)
}