package models

import (
	"strings"
	"testing"
)

func TestPost_RenderedBody(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		contains string
	}{
		{
			name:     "renders heading",
			body:     "## Hello World",
			contains: "<h2>Hello World</h2>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{Body: tt.body}
			result := string(post.RenderedBody())

			if tt.contains == "" && strings.TrimSpace(result) != "" {
				t.Errorf("expected empty output, go %q", result)
				return
			}

			if !strings.Contains(result, tt.contains) {
				t.Errorf("expected output to contain %q\ngot: %s", tt.contains, result)
			}
		})
	}
}
