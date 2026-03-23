package queries

import (
	"context"
	"fmt"

	"github.com/hiimtaylorjones/hiimtaylor-go/database"
	"github.com/hiimtaylorjones/hiimtaylor-go/models"
)

func CountPublishedPosts() (int, error) {
	var count int
	err := database.Pool.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM posts WHERE published = TRUE`,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting posts: %w", err)
	}
	return count, nil
}

func GetPublishedPosts(page, perPage int) ([]models.Post, error) {
	offset := (page - 1) * perPage
	query := `SELECT id, title, tagline, body, slug, banner_image_url, published, created_at, updated_at 
						FROM posts WHERE published = TRUE 
						ORDER BY created_at DESC
						LIMIT $1 OFFSET $2`

	rows, err := database.Pool.Query(
		context.Background(),
		query,
		perPage, offset,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
			&p.BannerImageURL, &p.Published, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error parsing post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func GetPostBySlug(slug string) (models.Post, error) {
	var p models.Post
	query := `SELECT id, title, tagline, body, slug, published, banner_image_url, created_at, updated_at 
						FROM posts WHERE slug = $1`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		slug,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.BannerImageURL, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return models.Post{}, fmt.Errorf("post not found: %w", err)
	}
	return p, nil
}

func CreatePost(title, tagline, body, slug, bannerImageURL string, published bool) (models.Post, error) {
	var p models.Post
	query := `
		INSERT INTO posts (title, tagline, body, slug, published, banner_image_url) 
			VALUES($1, $2, $3, $4, $5, $6) 
			RETURNING id, title, tagline, body, slug, published, banner_image_url, created_at, updated_at
	`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		title, tagline, body, slug, published, bannerImageURL,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.BannerImageURL, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return models.Post{}, fmt.Errorf("error creating post: %w", err)
	}
	return p, nil
}

func UpdatePost(id int, title, tagline, body, bannerImageURL string, published bool) (models.Post, error) {
	var p models.Post
	query := `
		UPDATE posts SET title=$1, tagline=$2, body=$3, published=$4, banner_image_url=$5, updated_at=NOW()
			WHERE id=$6
			RETURNING id, title, tagline, body, slug, published, banner_image_url, created_at, updated_at
	`

	err := database.Pool.QueryRow(
		context.Background(),
		query,
		title, tagline, body, published, bannerImageURL, id,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.BannerImageURL, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return models.Post{}, fmt.Errorf("error updating post: %w", err)
	}
	return p, nil
}

func DeletePost(id int) error {
	_, err := database.Pool.Exec(
		context.Background(),
		"DELETE FROM posts WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}

func GetAdminByEmail(email string) (models.Admin, error) {
	var a models.Admin
	err := database.Pool.QueryRow(
		context.Background(),
		`SELECT id, email, encrypted_password FROM admins WHERE email = $1`,
		email,
	).Scan(&a.ID, &a.Email, &a.EncryptedPassword)

	if err != nil {
		return models.Admin{}, fmt.Errorf("admin not found: %w", err)
	}
	return a, nil
}

func CreateAdmin(email, hashedPassword string) (models.Admin, error) {
	var a models.Admin
	err := database.Pool.QueryRow(
		context.Background(),
		`INSERT INTO admins (email, encrypted_password)
						VALUES ($1, $2)
						RETURNING id, email, encrypted_password`,
		email, hashedPassword,
	).Scan(&a.ID, &a.Email, &a.EncryptedPassword)
	if err != nil {
		return models.Admin{}, fmt.Errorf("error creating admin: %w", err)
	}
	return a, nil
}
