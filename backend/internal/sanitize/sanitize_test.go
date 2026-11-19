package sanitize_test

import (
	"strings"
	"testing"

	"github.com/h00s/newsfuse/internal/sanitize"
)

func TestStoryRemovesScripts(t *testing.T) {
	got := sanitize.Story(`<p>Tekst</p><script>alert(1)</script><p onclick="x()">Drugi</p>`)
	if strings.Contains(got, "script") || strings.Contains(got, "onclick") {
		t.Errorf("Story kept active content: %q", got)
	}
	if !strings.Contains(got, "<p>Tekst</p>") || !strings.Contains(got, "<p>Drugi</p>") {
		t.Errorf("Story dropped paragraphs: %q", got)
	}
}

func TestStoryOpensLinksSafely(t *testing.T) {
	got := sanitize.Story(`<p><a href="https://example.test/a">link</a></p>`)
	for _, want := range []string{`href="https://example.test/a"`, `target="_blank"`, "nofollow", "noreferrer"} {
		if !strings.Contains(got, want) {
			t.Errorf("Story(link) = %q, missing %s", got, want)
		}
	}
}

func TestStoryDropsJavascriptURLs(t *testing.T) {
	got := sanitize.Story(`<p><a href="javascript:alert(1)">link</a></p>`)
	if strings.Contains(got, "javascript") {
		t.Errorf("Story kept a javascript: URL: %q", got)
	}
}

func TestPlainText(t *testing.T) {
	got := sanitize.PlainText(`<p>Prvi &amp; drugi</p><p>Treći</p>`)
	if strings.ContainsAny(got, "<>") || strings.Contains(got, "&amp;") {
		t.Errorf("PlainText left markup or entities: %q", got)
	}
	if !strings.Contains(got, "Prvi & drugi") || !strings.Contains(got, "Treći") {
		t.Errorf("PlainText lost text: %q", got)
	}
}

func TestStoryKeepsImages(t *testing.T) {
	got := sanitize.Story(`<p><img decoding="async" class="aligncenter" src="https://example.test/a.jpg" alt="Slika" width="640" height="480" srcset="https://example.test/a-300.jpg 300w"></p>`)
	for _, want := range []string{`src="https://example.test/a.jpg"`, `alt="Slika"`, `width="640"`, `height="480"`} {
		if !strings.Contains(got, want) {
			t.Errorf("Story(img) = %q, missing %s", got, want)
		}
	}
	for _, unwanted := range []string{"class=", "srcset", "decoding"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("Story(img) = %q, kept %s", got, unwanted)
		}
	}
}

func TestStoryDropsUnsafeImages(t *testing.T) {
	got := sanitize.Story(`<p><img src="javascript:alert(1)"><img src="data:image/svg+xml,%3Csvg%3E" alt="x"><img src="https://example.test/b.jpg" onerror="steal()"></p>`)
	if strings.Contains(got, "javascript") || strings.Contains(got, "data:") || strings.Contains(got, "onerror") {
		t.Errorf("Story kept an unsafe image: %q", got)
	}
	if strings.Count(got, "<img") != 1 {
		t.Errorf("Story = %q, want only the https image, and no image left without a source", got)
	}
}
