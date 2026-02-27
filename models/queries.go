package models

import (
	"context"
	"fmt"

	"github.com/hiimtaylorjones/hiimtaylor-go/database"
)

func GetPublishedPosts() ([]Post, error) {
	query := "SELECT id, title, tagline, body, slug, published, created_at, updated_at FROM posts WHERE published = TRUE ORDER BY created_at DESC"

	rows, err := database.Pool.Query(
		context.Background(), query,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
			&p.Published, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error parsing post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func GetPostBySlug(slug string) (Post, error) {
	var p Post
	query := "SELECT id, title, tagline, body, slug, published, created_at, updated_at FROM posts WHERE slug = $1"
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		slug,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		return Post{}, fmt.Errorf("post not found: %w", err)
	}
	return p, nil
}

func CreatePost(title, tagline, body, slug string, published bool) (Post, error) {
	var p Post
	query := `
		INSERT INTO posts (title, tagline, body, slug, published) 
			VALUES($1, $2, $3, $4, $5) 
			RETURNING id, tagline, body, slug, published, created_at, updated_at
	`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		title, tagline, body, slug, published,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.CreatedAt, &p.UpdatedAt,
	)
	
	if err != nil {
		return Post{}, fmt.Errorf("error creating post: %w", err)
	}
	return p, nil
}

func UpdatePost(id int, title, tagline, body string, published bool) (Post, error) {
	var p Post
	query := `
		UPDATE posts SET title=$1, tagline=$2, body=$3, published=$4, updated_at=NOW()
			WHERE id=$5
			RETURNING id, tagline, body, slug, published, created_at, updated_at  
	`

	err := database.Pool.QueryRow(
		context.Background(),
		query,
		title, tagline, body, published, id,
	).Scan(
		&p.ID, &p.Title, &p.Tagline, &p.Body, &p.Slug,
		&p.Published, &p.CreatedAt, &p.UpdatedAt,
	)
	
	if err != nil {
		return Post{}, fmt.Errorf("error updating post: %w", err)
	}
	return p, nil
}

func DeletePost(id int) error {
	_, err := database.Pool.Exec(
		context.Background(),
		"DELETE FROM post WHERE id=$1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}