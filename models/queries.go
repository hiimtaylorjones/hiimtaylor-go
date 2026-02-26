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
