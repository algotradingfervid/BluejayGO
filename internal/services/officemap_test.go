package services

import (
	"net/url"
	"strings"
	"testing"
)

func TestOfficeMapURLs(t *testing.T) {
	address := "123 Main Street, Hyderabad, India"
	for _, raw := range []string{"", "https://maps.app.goo.gl/Example", "https://goo.gl/maps/Example", "https://www.google.com/maps/place/BlueJay+Office/", "https://maps.google.com/?q=Hyderabad", "https://www.google.com/maps/embed?pb=!1m18!2m3"} {
		if !ValidOfficeMapURL(raw) {
			t.Fatalf("valid map rejected: %s", raw)
		}
		embed, directions := OfficeMapURLs(raw, address)
		u, err := url.Parse(embed)
		if err != nil || u.Scheme != "https" || u.Host != "www.google.com" || !strings.HasPrefix(u.Path, "/maps") {
			t.Fatalf("unsafe frame %s", embed)
		}
		if raw != "" && directions != raw {
			t.Fatalf("directions lost stored URL %s", raw)
		}
		if strings.Contains(raw, "goo.gl") && u.Query().Get("q") != address {
			t.Fatal("short link must fall back to address without outbound request")
		}
	}
	for _, raw := range []string{"javascript:alert(1)", "https://evil.example/map", "http://maps.google.com/", "https://www.google.com.evil.example/maps", "https://user@www.google.com/maps", "https://www.google.com:123/maps", "https://www.google.com/url?url=evil"} {
		if ValidOfficeMapURL(raw) {
			t.Fatalf("unsafe map accepted: %s", raw)
		}
		embed, directions := OfficeMapURLs(raw, address)
		if strings.Contains(embed, "evil") || strings.Contains(directions, "evil") {
			t.Fatal("unsafe stored URL must fall back to address")
		}
	}
}
