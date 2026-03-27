package queries

import (
	"log"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/hiimtaylorjones/hiimtaylor-go/database"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	godotenv.Load()
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

func TestUpdatePost(t *testing.T) {
	post, err := CreatePost("Test Post", "tag", "body", "test-post", "", false)
	if err != nil {
		t.Fatalf("Error creating post: %v", err)
	}

	post, err = UpdatePost(post.ID, "Test Post", "tag", "body", "", true)

	if err != nil {
		t.Fatalf("Error updating post: %v", err)
	}

	if post.Published != true {
		t.Errorf("expected post update to make published true")
	}

	t.Cleanup(func() {
		DeletePost(post.ID)
	})
}

func TestAdminCreateFetchFlow(t *testing.T) {
	var password string = "my-secret-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	admin, err := CreateAdmin("hello!@me.com", string(hash))

	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}

	if admin.Email == "" {
		t.Error("expected a non-zero ID from Return")
	}

	admin, err = GetAdminByEmail(admin.Email)

	if err != nil {
		t.Fatalf("Error fetching admin: %v", err)
	}

	t.Cleanup(func() {
		DeleteAdmin(admin.Email)
	})
}
