package content

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/yuin/goldmark"
)

func Render(filename string) (template.HTML, error) {
	raw, err := os.ReadFile("content/" + filename)
	if err != nil {
		return "", fmt.Errorf("could not read content file: %w", err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(raw, &buf); err != nil {
		return "", fmt.Errorf("could not render markdown: %w", err)
	}

	return template.HTML(buf.String()), nil
}