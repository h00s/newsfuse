package utils_test

import (
	"strings"
	"testing"

	"github.com/h00s/newsfuse/app/utils"
)

func TestSanitizeStoryRemovesScripts(t *testing.T) {
	got := utils.SanitizeStory(`<p>Tekst</p><script>alert(1)</script><p onclick="x()">Drugi</p>`)
	if strings.Contains(got, "script") || strings.Contains(got, "onclick") {
		t.Errorf("SanitizeStory kept active content: %q", got)
	}
	if !strings.Contains(got, "<p>Tekst</p>") || !strings.Contains(got, "<p>Drugi</p>") {
		t.Errorf("SanitizeStory dropped paragraphs: %q", got)
	}
}

func TestSanitizeStoryOpensLinksSafely(t *testing.T) {
	got := utils.SanitizeStory(`<p><a href="https://example.test/a">link</a></p>`)
	for _, want := range []string{`href="https://example.test/a"`, `target="_blank"`, "nofollow", "noreferrer"} {
		if !strings.Contains(got, want) {
			t.Errorf("SanitizeStory(link) = %q, missing %s", got, want)
		}
	}
}

func TestSanitizeStoryDropsJavascriptURLs(t *testing.T) {
	got := utils.SanitizeStory(`<p><a href="javascript:alert(1)">link</a></p>`)
	if strings.Contains(got, "javascript") {
		t.Errorf("SanitizeStory kept a javascript: URL: %q", got)
	}
}

func TestStoryTextIsPlainText(t *testing.T) {
	got := utils.StoryText(`<p>Prvi &amp; drugi</p><p>Treći</p>`)
	if strings.ContainsAny(got, "<>") || strings.Contains(got, "&amp;") {
		t.Errorf("StoryText left markup or entities: %q", got)
	}
	if !strings.Contains(got, "Prvi & drugi") || !strings.Contains(got, "Treći") {
		t.Errorf("StoryText lost text: %q", got)
	}
}

func TestSanitizeStoryKeepsImages(t *testing.T) {
	got := utils.SanitizeStory(`<p><img decoding="async" class="aligncenter" src="https://example.test/a.jpg" alt="Slika" width="640" height="480" srcset="https://example.test/a-300.jpg 300w"></p>`)
	for _, want := range []string{`src="https://example.test/a.jpg"`, `alt="Slika"`, `width="640"`, `height="480"`} {
		if !strings.Contains(got, want) {
			t.Errorf("SanitizeStory(img) = %q, missing %s", got, want)
		}
	}
	for _, unwanted := range []string{"class=", "srcset", "decoding"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("SanitizeStory(img) = %q, kept %s", got, unwanted)
		}
	}
}

func TestSanitizeStoryDropsUnsafeImages(t *testing.T) {
	got := utils.SanitizeStory(`<p><img src="javascript:alert(1)"><img src="data:image/svg+xml,%3Csvg%3E" alt="x"><img src="https://example.test/b.jpg" onerror="steal()"></p>`)
	if strings.Contains(got, "javascript") || strings.Contains(got, "data:") || strings.Contains(got, "onerror") {
		t.Errorf("SanitizeStory kept an unsafe image: %q", got)
	}
	if strings.Count(got, "<img") != 1 {
		t.Errorf("SanitizeStory = %q, want only the https image, and no image left without a source", got)
	}
}
