package slug

import (
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9-]+`)

func Generate(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphanumeric.ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}