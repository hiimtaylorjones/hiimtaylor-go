package queries

import (
	"log"
	"os"
	"testing"

	"github.com/hiimtaylorjones/hiimtaylor-go/database"
)

func TestMain(m *testing.M) {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL must be set to run tests")
	}
	database.Connect()
	defer database.Close()
	os.Exit(m.Run())
}

func TestCreatePost(t *testing.T) {
	post, err := CreatePost("Queries Test", "tagline", "body", "queries-test", "", false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	t.Cleanup(func() {
		DeletePost(post.ID)
	})

	if post.ID == 0 {
		t.Error("expected a non-zero ID from Return")
	}
	if post.Title != "Queries Test" {
		t.Errorf("expected title %q, got %q", "Queries Test", post.Title)
	}
}
