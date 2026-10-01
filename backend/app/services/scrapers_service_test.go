package services

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
)

// fakeScraper counts scrapes and never touches the network. Its zero schedule waits nothing
// between scrapes, so a loop's timer is always ready when its context is cancelled.
type fakeScraper struct{ scrapes atomic.Int32 }

func (f *fakeScraper) String() string           { return "fake" }
func (f *fakeScraper) Schedule() utils.Schedule { return utils.Schedule{} }

func (f *fakeScraper) Scrape(ctx context.Context) (models.Headlines, error) {
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
	scraper := &fakeScraper{}
	shutdown()

	done := make(chan struct{})
	go func() {
		s.run(scraper)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the loop kept running after shutdown")
	}
	if n := scraper.scrapes.Load(); n != 0 {
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
