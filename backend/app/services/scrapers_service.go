package services

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/go-raptor/raptor/v4/errs"
	"github.com/h00s/newsfuse/app/utils"
	"github.com/h00s/newsfuse/app/utils/scrapers"
)

// ScrapersService runs one loop per news site, storing new headlines through HeadlinesService,
// and scrapes a story on demand. Each site is registered under the id of its source row.
type ScrapersService struct {
	raptor.Service

	Headlines *HeadlinesService

	scrapers map[int64]utils.Scraper
	cancel   context.CancelFunc
	running  sync.WaitGroup
}

func (s *ScrapersService) Setup() error {
	s.scrapers = map[int64]utils.Scraper{
		1:  scrapers.NewKliknihr(1),
		2:  scrapers.NewMojportalhr(2),
		3:  scrapers.NewRadioDaruvar(3),
		4:  scrapers.NewIndexhrCroatia(4),
		5:  scrapers.NewN1InfoCroatia(5),
		6:  scrapers.NewIndexhrWorld(6),
		7:  scrapers.NewN1InfoWorld(7),
		8:  scrapers.NewHackerNews(8),
		9:  scrapers.NewBughr(9),
		10: scrapers.NewTelegram(10),
		11: scrapers.NewHCL(11),
	}

	if s.Config.AppConfig["scrapers_enabled"] == "false" {
		s.Log.Info("Scrapers are disabled (scrapers_enabled: false)")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	for _, scraper := range s.scrapers {
		s.running.Go(func() { s.run(ctx, scraper) })
	}
	return nil
}

// Cleanup stops every loop and waits for them. A scrape in flight is canceled through its
// context; an ingest in flight finishes first, because the database outlives this service.
func (s *ScrapersService) Cleanup() error {
	if s.cancel != nil {
		s.cancel()
	}
	s.running.Wait()
	return nil
}

func (s *ScrapersService) run(ctx context.Context, scraper utils.Scraper) {
	schedule := scraper.Schedule()
	for {
		if hour := time.Now().Hour(); slices.Contains(schedule.OffHours, hour) {
			s.Log.Debug("Skipping scrape in off-hours", "scraper", scraper.String(), "hour", hour)
		} else {
			s.scrape(ctx, scraper)
		}

		wait := schedule.NextWait()
		s.Log.Debug("Waiting for the next scrape", "scraper", scraper.String(), "wait", wait.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}

func (s *ScrapersService) scrape(ctx context.Context, scraper utils.Scraper) {
	headlines, err := scraper.Scrape(ctx)
	if err != nil {
		if ctx.Err() == nil {
			s.Log.Warn("Scrape failed", "scraper", scraper.String(), "error", err)
		}
		return
	}
	inserted, err := s.Headlines.Ingest(headlines)
	if err != nil {
		s.Log.Error("Could not store headlines", "scraper", scraper.String(), "error", err)
		return
	}
	s.Log.Info("Scraped", "scraper", scraper.String(), "headlines", len(headlines), "new", inserted)
}

// ScrapeStory fetches the story behind a headline of the given source, as sanitized HTML.
func (s *ScrapersService) ScrapeStory(ctx context.Context, sourceID int64, url string) (string, error) {
	scraper, ok := s.scrapers[sourceID]
	if !ok {
		return "", errs.NewErrorNotFound("Story not available")
	}
	content, err := scraper.ScrapeStory(ctx, url)
	if err != nil {
		s.Log.Warn("Story scrape failed", "scraper", scraper.String(), "url", url, "error", err)
		return "", errs.NewErrorBadGateway("Could not fetch the story")
	}
	if strings.TrimSpace(content) == "" {
		s.Log.Warn("Story page had no text", "scraper", scraper.String(), "url", url)
		return "", errs.NewErrorBadGateway("Could not fetch the story")
	}
	return content, nil
}
