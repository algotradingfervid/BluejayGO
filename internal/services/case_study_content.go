package services

import (
	"encoding/json"
	"html"
	"strings"
	"unicode/utf8"
)

// NormalizeCaseStudySection preserves legacy prose entered into a heading field,
// presenting it as body copy under a short, meaningful section heading.
func NormalizeCaseStudySection(title, content, fallback string) (string, string) {
	title = strings.TrimSpace(title)
	if utf8.RuneCountInString(title) > 120 || strings.ContainsAny(title, "\r\n") || strings.HasPrefix(title, "•") {
		var body strings.Builder
		if strings.HasPrefix(title, "•") {
			body.WriteString("<ul>")
			for _, item := range strings.Split(title, "•") {
				if item = strings.TrimSpace(item); item != "" {
					body.WriteString("<li>" + html.EscapeString(item) + "</li>")
				}
			}
			body.WriteString("</ul>")
		} else {
			for _, line := range strings.Split(strings.ReplaceAll(title, "\r\n", "\n"), "\n") {
				if line = strings.TrimSpace(line); line != "" {
					body.WriteString("<p>" + html.EscapeString(line) + "</p>\n")
				}
			}
		}
		content = body.String() + "\n" + content
		title = fallback
	}
	if title == "" {
		title = fallback
	}
	return title, content
}

// ParseCaseStudyBullets accepts stored JSON and the editor's one-item-per-line
// format. Commas belong to an item, not to the item separator.
func ParseCaseStudyBullets(raw string) []string {
	var bullets []string
	if json.Unmarshal([]byte(raw), &bullets) != nil {
		bullets = strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n")
	} else {
		// Older edits displayed JSON in the textarea and split it at commas on save.
		// Repair that representation for display without discarding the stored text.
		var repaired []string
		if json.Unmarshal([]byte(strings.Join(bullets, ",")), &repaired) == nil {
			bullets = repaired
		}
	}
	cleaned := make([]string, 0, len(bullets))
	for _, bullet := range bullets {
		if bullet = strings.TrimSpace(bullet); bullet != "" {
			cleaned = append(cleaned, bullet)
		}
	}
	return cleaned
}
