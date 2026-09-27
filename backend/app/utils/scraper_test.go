package utils_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

func servePage(t *testing.T, body string) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func TestScrapeStoryEscapesTextParagraphs(t *testing.T) {
	url := servePage(t, `<html><body><div class="article">
		<p>Prvi &lt;script&gt;alert(1)&lt;/script&gt;</p>
		<p>   </p>
		<p>Drugi</p>
	</div></body></html>`)
	s := utils.NewScraper("test", url, 1, 2, nil)

	got, err := s.ScrapeStoryFrom(t.Context(), url, "div.article", "p", false)
	if err != nil {
		t.Fatal(err)
	}
	want := "<p>Prvi &lt;script&gt;alert(1)&lt;/script&gt;</p><p>Drugi</p>"
	if got != want {
		t.Errorf("ScrapeStoryFrom(text) = %q, want %q", got, want)
	}
}

func TestScrapeStorySanitizesHTMLParagraphs(t *testing.T) {
	url := servePage(t, `<html><body><table><tr><td class="title"><span class="titleline">
		<a href="https://example.test/a" onclick="steal()">Naslov</a><script>steal()</script>
	</span></td></tr></table></body></html>`)
	s := utils.NewScraper("test", url, 1, 2, nil)

	got, err := s.ScrapeStoryFrom(t.Context(), url, "td.title", "span.titleline", true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "script") || strings.Contains(got, "onclick") {
		t.Errorf("ScrapeStoryFrom(html) kept active content: %q", got)
	}
	if !strings.Contains(got, `href="https://example.test/a"`) || !strings.Contains(got, "Naslov") {
		t.Errorf("ScrapeStoryFrom(html) lost the link: %q", got)
	}
}

func TestScrapeCollectsHeadlines(t *testing.T) {
	url := servePage(t, `<html><body>
		<h3><a href="https://example.test/1"> Prvi naslov </a></h3>
		<h3><a href="https://example.test/2">   </a></h3>
		<h3><a href="https://example.test/3">Treći naslov</a></h3>
	</body></html>`)
	s := utils.NewScraper("test", url, 1, 2, nil)
	s.ScrapeHeadline("h3", func(e *colly.HTMLElement) {
		s.AddHeadline(models.Headline{SourceID: 7, Title: e.ChildText("a"), URL: e.ChildAttr("a", "href")})
	})

	got, err := s.Scrape(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Title != "Prvi naslov" || got[1].URL != "https://example.test/3" || got[0].SourceID != 7 {
		t.Errorf("Scrape = %+v, want the two non-empty headlines in page order", got)
	}

	again, err := s.Scrape(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 2 {
		t.Errorf("second Scrape returned %d headlines, want 2 (each run starts empty)", len(again))
	}
}

func TestScrapeStopsWhenCanceled(t *testing.T) {
	url := servePage(t, `<html><body><h3>Naslov</h3></body></html>`)
	s := utils.NewScraper("test", url, 1, 2, nil)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, err := s.Scrape(ctx); err == nil {
		t.Error("Scrape with a canceled context succeeded, want an error")
	}
}

func TestScheduleNextWait(t *testing.T) {
	fixed := utils.NewScraper("test", "https://example.test", 15, 15, nil).Schedule()
	if got := fixed.NextWait(); got != 15*time.Minute {
		t.Errorf("NextWait with min == max = %v, want 15m", got)
	}

	ranged := utils.NewScraper("test", "https://example.test", 10, 20, []int{3}).Schedule()
	for range 100 {
		if got := ranged.NextWait(); got < 10*time.Minute || got >= 20*time.Minute {
			t.Fatalf("NextWait = %v, want within [10m, 20m)", got)
		}
	}
	if len(ranged.OffHours) != 1 || ranged.OffHours[0] != 3 {
		t.Errorf("OffHours = %v, want [3]", ranged.OffHours)
	}
}
