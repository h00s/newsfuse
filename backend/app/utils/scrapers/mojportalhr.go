package scrapers

import (
	"context"
	"time"

	"github.com/h00s/newsfuse/app/models"
	"github.com/h00s/newsfuse/app/utils"
	"github.com/h00s/newsfuse/internal/collector"
)

// mojportalURL is the site's canonical host; www.mojportal.hr redirects here.
const mojportalURL = "https://mojportal.hr"

type Mojportalhr struct {
	utils.DefaultScraper
}

func NewMojportalhr(sourceID int64) *Mojportalhr {
	s := &Mojportalhr{
		DefaultScraper: *utils.NewScraper(
			"MojPortal.hr",
			mojportalURL+"/",
			15,
			30,
			[]int{0, 1, 2, 3, 4, 5, 23},
		),
	}

	// News cards only: the front page also lists event-calendar cards (se-card--event) and
	// classifieds (se-card--classified). Links are site-relative.
	s.ScrapeHeadline("article.se-card--news:not(.se-card--event) > a.se-card__wrap", func(e *collector.HTMLElement) {
		s.AddHeadline(models.Headline{
			SourceID:    sourceID,
			Title:       e.ChildText(".se-card__head"),
			URL:         mojportalURL + e.Attr("href"),
			PublishedAt: time.Now(),
		})
	})

	return s
}

// ScrapeStory skips the closing "follow us on Facebook and TikTok" paragraph.
func (s *Mojportalhr) ScrapeStory(ctx context.Context, url string) (string, error) {
	return s.ScrapeStoryFrom(ctx, url, "div.se-article__text", "p:not([class]):not(:has(a[href*='facebook.com/MojPortal']))", false)
}
