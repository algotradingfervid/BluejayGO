package services

import (
	"net/url"
	"strings"
)

// ValidOfficeMapURL accepts Google Maps URLs, including common share links.
// Short links are used for directions; the embedded map falls back to the address.
// We never fetch an administrator-supplied URL from the server.
func ValidOfficeMapURL(raw string) bool {
	if raw == "" {
		return true
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return false
	}
	switch strings.ToLower(u.Hostname()) {
	case "maps.app.goo.gl":
		return u.Path != "" && u.Path != "/"
	case "goo.gl":
		return strings.HasPrefix(u.Path, "/maps/")
	case "maps.google.com", "maps.google.co.in":
		return true
	case "google.com", "www.google.com", "google.co.in", "www.google.co.in":
		return u.Path == "/maps" || strings.HasPrefix(u.Path, "/maps/")
	}
	return false
}

// OfficeMapURLs always returns an embed on the single CSP-allowed Google origin.
func OfficeMapURLs(raw, address string) (embed, directions string) {
	query := address
	directions = "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(address)
	if raw != "" && ValidOfficeMapURL(raw) {
		directions = raw
		u, _ := url.Parse(raw)
		host := strings.ToLower(u.Hostname())
		if host != "maps.app.goo.gl" && host != "goo.gl" {
			if strings.HasPrefix(u.Path, "/maps/embed") && u.Query().Get("pb") != "" {
				return "https://www.google.com/maps/embed?pb=" + url.QueryEscape(u.Query().Get("pb")), directions
			}
			if q := u.Query().Get("q"); q != "" {
				query = q
			} else if q := u.Query().Get("query"); q != "" {
				query = q
			} else if q := u.Query().Get("destination"); q != "" {
				query = q
			} else if strings.HasPrefix(u.Path, "/maps/place/") {
				place := strings.Split(strings.TrimPrefix(u.Path, "/maps/place/"), "/")[0]
				if place != "" {
					query = strings.ReplaceAll(place, "+", " ")
				}
			}
		}
	}
	return "https://www.google.com/maps?q=" + url.QueryEscape(query) + "&output=embed", directions
}
