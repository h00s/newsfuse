package services

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/internal/scraper"
	"github.com/h00s/newsfuse/internal/scraper/sites"
)

// ScrapersService runs one loop per news site, storing new headlines through HeadlinesService,
// and scrapes a story on demand. Each site is registered under the id of its source row.
type ScrapersService struct {
	raptor.Service

	Headlines *HeadlinesService

	scrapers map[int64]scraper.Scraper
	running  sync.WaitGroup
}

func (s *ScrapersService) Setup() error {
	s.scrapers = map[int64]scraper.Scraper{
		1:  sites.NewKliknihr(),
		2:  sites.NewMojportalhr(),
		3:  sites.NewRadioDaruvar(),
		4:  sites.NewIndexhrCroatia(),
		5:  sites.NewN1InfoCroatia(),
		6:  sites.NewIndexhrWorld(),
		7:  sites.NewN1InfoWorld(),
		8:  sites.NewHackerNews(),
		9:  sites.NewBughr(),
		10: sites.NewTelegram(),
		11: sites.NewHCL(),
	}

	enabled, err := s.Config.AppBool("scrapers_enabled", true)
	if err != nil {
		return err
	}
	if !enabled {
		s.Log.Info("Scrapers are disabled (scrapers_enabled: false)")
		return nil
	}

	for id, sc := range s.scrapers {
		s.running.Go(func() { s.run(id, sc) })
	}
	return nil
}

// Cleanup waits for the loops, which return once the app context is cancelled; a scrape or an
// ingest in flight is cancelled with it. Waiting matters because an ingest invalidates the
// cache, which is cleaned up after this service.
func (s *ScrapersService) Cleanup() error {
	s.running.Wait()
	return nil
}

// run scrapes on the site's schedule until the app context is cancelled at shutdown.
func (s *ScrapersService) run(sourceID int64, sc scraper.Scraper) {
	ctx := s.AppContext()
	schedule := sc.Schedule()
	for {
		// Checked first: when the timer and the cancellation are both ready, select picks either,
		// and no scrape may start after shutdown.
		if ctx.Err() != nil {
			return
		}
		if hour := time.Now().Hour(); slices.Contains(schedule.OffHours, hour) {
			s.Log.Debug("Skipping scrape in off-hours", "scraper", sc.String(), "hour", hour)
		} else {
			s.scrape(ctx, sourceID, sc)
		}

		wait := schedule.NextWait()
		s.Log.Debug("Waiting for the next scrape", "scraper", sc.String(), "wait", wait.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}

func (s *ScrapersService) scrape(ctx context.Context, sourceID int64, sc scraper.Scraper) {
	scraped, err := sc.Scrape(ctx)
	if err != nil {
		if ctx.Err() == nil {
			s.Log.Warn("Scrape failed", "scraper", sc.String(), "error", err)
		}
		return
	}
	headlines := toHeadlines(sourceID, scraped, time.Now())
	inserted, err := s.Headlines.Ingest(headlines)
	if err != nil {
		// Shutdown cancels an ingest in flight, and the next boot scrapes those headlines again.
		if ctx.Err() == nil {
			s.Log.Error("Could not store headlines", "scraper", sc.String(), "error", err)
		}
		return
	}
	s.Log.Info("Scraped", "scraper", sc.String(), "headlines", len(headlines), "new", inserted)
}

// toHeadlines files a site's scraped headlines under its source. One timestamp serves the whole
// batch: lists order headlines by id, and Ingest gives newer headlines higher ids.
func toHeadlines(sourceID int64, scraped []scraper.Headline, now time.Time) models.Headlines {
	headlines := make(models.Headlines, len(scraped))
	for i, h := range scraped {
		headlines[i] = models.Headline{SourceID: sourceID, Title: h.Title, URL: h.URL, PublishedAt: now}
	}
	return headlines
}

// ScrapeStory fetches the story behind a headline of the given source, as sanitized HTML.
func (s *ScrapersService) ScrapeStory(ctx context.Context, sourceID int64, url string) (string, error) {
	sc, ok := s.scrapers[sourceID]
	if !ok {
		return "", errs.NewErrorNotFound("Story not available")
	}
	content, err := sc.ScrapeStory(ctx, url)
	if err != nil {
		s.Log.Warn("Story scrape failed", "scraper", sc.String(), "url", url, "error", err)
		return "", errs.NewErrorBadGateway("Could not fetch the story")
	}
	if strings.TrimSpace(content) == "" {
		s.Log.Warn("Story page had no text", "scraper", sc.String(), "url", url)
		return "", errs.NewErrorBadGateway("Could not fetch the story")
	}
	return content, nil
}
