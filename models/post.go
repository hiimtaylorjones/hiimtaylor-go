package models

import "time"

type Post struct {
	ID 				int
	Title 		string
	Tagline		string
	Body			string
	Slug			string
	Published	bool
	CreatedAt	time.Time
	UpdatedAt	time.Time
}