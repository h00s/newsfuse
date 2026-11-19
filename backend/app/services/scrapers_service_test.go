package services

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal/scraper"
)

// fakeScraper counts scrapes and never touches the network. Its zero schedule waits nothing
// between scrapes, so a loop's timer is always ready when its context is cancelled.
type fakeScraper struct{ scrapes atomic.Int32 }

func (f *fakeScraper) String() string             { return "fake" }
func (f *fakeScraper) Schedule() scraper.Schedule { return scraper.Schedule{} }

func (f *fakeScraper) Scrape(ctx context.Context) ([]scraper.Headline, error) {
	f.scrapes.Add(1)
	return nil, ctx.Err()
}

func (f *fakeScraper) ScrapeStory(context.Context, string) (string, error) { return "", nil }

func TestScraperLoopStopsWithTheApp(t *testing.T) {
	res, shutdown := testResources()
	s := &ScrapersService{}
	if err := s.Init(res); err != nil {
		t.Fatal(err)
	}
	fake := &fakeScraper{}
	shutdown()

	done := make(chan struct{})
	go func() {
		s.run(0, fake)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the loop kept running after shutdown")
	}
	if n := fake.scrapes.Load(); n != 0 {
		t.Errorf("scraped %d times after shutdown, want 0", n)
	}
}

func TestScrapersRejectAnUnparseableSwitch(t *testing.T) {
	res, shutdown := testResources()
	res.Config.AppConfig = map[string]string{"scrapers_enabled": "nope"}
	shutdown() // should the loops start anyway, they return before scraping
	s := &ScrapersService{}
	if err := s.Init(res); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Cleanup() })

	err := s.Setup()
	if err == nil || !strings.Contains(err.Error(), "scrapers_enabled") {
		t.Errorf("Setup = %v, want an error naming scrapers_enabled", err)
	}
}

func TestToHeadlinesStampsSourceAndTime(t *testing.T) {
	now := time.Date(2026, 11, 19, 12, 0, 0, 0, time.UTC)
	got := toHeadlines(4, []scraper.Headline{{Title: "A", URL: "https://x/a"}, {Title: "B", URL: "https://x/b"}}, now)

	want := models.Headlines{
		{SourceID: 4, Title: "A", URL: "https://x/a", PublishedAt: now},
		{SourceID: 4, Title: "B", URL: "https://x/b", PublishedAt: now},
	}
	if len(got) != len(want) {
		t.Fatalf("toHeadlines returned %d headlines, want %d", len(got), len(want))
	}
	for i, w := range want {
		g := got[i]
		if g.SourceID != w.SourceID || g.Title != w.Title || g.URL != w.URL || !g.PublishedAt.Equal(w.PublishedAt) {
			t.Errorf("headline %d = {%d %q %q %v}, want {%d %q %q %v}", i, g.SourceID, g.Title, g.URL, g.PublishedAt, w.SourceID, w.Title, w.URL, w.PublishedAt)
		}
	}
}
