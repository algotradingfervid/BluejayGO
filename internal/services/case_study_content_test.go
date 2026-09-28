package services

import (
	"reflect"
	"strings"
	"testing"
)

func TestNormalizeCaseStudySectionPreservesLegacyProse(t *testing.T) {
	prose := strings.Repeat("Community health staff require dependable attendance. ", 5)
	title, content := NormalizeCaseStudySection(prose, "<p>Existing body.</p>", "The Challenge")
	if title != "The Challenge" || !strings.Contains(content, strings.TrimSpace(prose)) || !strings.Contains(content, "<p>Existing body.</p>") {
		t.Fatalf("legacy content lost: %q %q", title, content)
	}
	title2, content2 := NormalizeCaseStudySection(title, content, "The Challenge")
	if title2 != title || content2 != content {
		t.Fatal("normalization duplicated content on repeat edit")
	}
	_, escaped := NormalizeCaseStudySection(strings.Repeat("<script>alert(1)</script>", 6), "", "The Challenge")
	if strings.Contains(escaped, "<script>") {
		t.Fatal("title text became executable markup")
	}
	title, content = NormalizeCaseStudySection("Field connectivity", "<p>Body</p>", "The Challenge")
	if title != "Field connectivity" || content != "<p>Body</p>" {
		t.Fatal("short heading changed")
	}
}

func TestCaseStudyBulletsPreserveCommasAndRepairLegacyJSON(t *testing.T) {
	want := []string{"GPS, Wi-Fi and Bluetooth", "Reliable attendance"}
	for _, raw := range []string{
		"GPS, Wi-Fi and Bluetooth\r\nReliable attendance\n",
		`["GPS, Wi-Fi and Bluetooth","Reliable attendance"]`,
		`["[\"GPS", "Wi-Fi and Bluetooth\"", "\"Reliable attendance\"]"]`,
	} {
		got := ParseCaseStudyBullets(raw)
		// A legacy comma split trimmed the spaces; it cannot reconstruct lost whitespace.
		if strings.Contains(raw, `Wi-Fi and Bluetooth\"`) {
			want = []string{"GPS,Wi-Fi and Bluetooth", "Reliable attendance"}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %#v want %#v", raw, got, want)
		}
	}
}

func TestNormalizeCaseStudyOutcomeBullets(t *testing.T) {
	title, content := NormalizeCaseStudySection("•1,704 terminals deployed •Daily attendance enabled •Certified devices", "", "The Outcome")
	if title != "The Outcome" || strings.Count(content, "<li>") != 3 || !strings.Contains(content, "<ul>") {
		t.Fatalf("outcome list not formatted: %q %q", title, content)
	}
}
