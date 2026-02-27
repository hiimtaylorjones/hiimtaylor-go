package slug

import (
	"regex"
	"strings"
)

var nonAlphanumeric = regex.MustCompile(`[^a-z0-9-]+`)

func Generate(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphanumeric.ReplaceAllStrings(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllStrings(s, "-")
	s = strings.Trim(s, "-")
}