package models

import (
	"bytes"
	"html/template"
	"time"

	"github.com/yuin/goldmark"
)

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

func (p Post) RenderedBody() template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(p.Body), &buf); err != nil {
		return template.HTML(p.Body)
	}
	return template.HTML(buf.String())
}