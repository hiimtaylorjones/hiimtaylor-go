package main

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/hiimtaylorjones/hiimtaylor-go/database"
)

func TestMain(m *testing.M) {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL must be set to run tests")
	}
	database.Connect()
	defer database.Close()

	sessionManager = scs.New()
	loadTemplates()

	os.Exit(m.Run())
}

func TestCreatePost_Success(t *testing.T) {
	body, contentType := buildPostForm(t, map[string]string{
		"title":     "Test Post",
		"tagline":   "A test tageline",
		"body":      "Hello",
		"published": "true",
	})

	req := httptest.NewRequest("POST", "/posts", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	handleCreatePost(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 SeeOther, got %d\nbody: %s", rr.Code, rr.Body.String())
	}

	location := rr.Header().Get("Location")
	if location == "" {
		t.Error("expected redirect location header, got none")
	}

	slug := strings.TrimPrefix(rr.Header().Get("Location"), "/posts/")
	t.Cleanup(func() { cleanupPostBySlug(t, slug) })
}

func TestCreatePost_MissingMultipartEncoding(t *testing.T) {
	// Simulates what happens when enctype="multitype/form-data" is missing
	req := httptest.NewRequest("POST", "/posts",
		bytes.NewBufferString("title=Test+Post&body=Hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handleCreatePost(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestCreatePost_UnpublishedByDefault(t *testing.T) {
	body, contentType := buildPostForm(t, map[string]string{
		"title": "Unpublished",
		"body":  "Draft",
	})

	req := httptest.NewRequest("POST", "/posts", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	handleCreatePost(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", rr.Code, rr.Body.String())
	}

	slug := strings.TrimPrefix(rr.Header().Get("Location"), "/posts/")
	t.Cleanup(func() { cleanupPostBySlug(t, slug) })
}

// Helpers

func buildPostForm(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for key, val := range fields {
		if err := w.WriteField(key, val); err != nil {
			t.Fatalf("error writing field %s: %v", key, err)
		}
	}
	w.Close()
	return &buf, w.FormDataContentType()
}

func cleanupPostBySlug(t *testing.T, slug string) {
	t.Helper()
	_, err := database.Pool.Exec(
		context.Background(),
		"DELETE FROM posts WHERE slug = $1",
		slug,
	)
	if err != nil {
		t.Logf("warning: could not clean up post %q: %v", slug, err)
	}
}
