// Package utils provides the scraper base that every news site builds on, and the HTML
// sanitizing that makes scraped stories safe to render.
package utils

import (
	"context"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/h00s/newsfuse/app/models"
)

const (
	requestTimeout     = 15 * time.Second
	headlinesUserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:151.0) Gecko/20100101 Firefox/151.0"
	storyUserAgent     = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
)

// Scraper is one news site: its front page yields headlines, and each headline's page yields its
// story. String names the site in logs.
type Scraper interface {
	String() string
	Schedule() Schedule
	Scrape(ctx context.Context) (models.Headlines, error)
	ScrapeStory(ctx context.Context, url string) (string, error)
}

// Schedule is how often a site is scraped: a random wait between MinInterval and MaxInterval,
// so requests don't arrive like clockwork, and no scraping during OffHours (local hours, 0-23).
type Schedule struct {
	MinInterval time.Duration
	MaxInterval time.Duration
	OffHours    []int
}

func (s Schedule) NextWait() time.Duration {
	if s.MaxInterval <= s.MinInterval {
		return s.MinInterval
	}
	return s.MinInterval + rand.N(s.MaxInterval-s.MinInterval)
}

type headlineRule struct {
	selector string
	callback colly.HTMLCallback
}

// DefaultScraper does the fetching; each site adds its headline rules with ScrapeHeadline and
// its story selectors in ScrapeStory. Scrape is not safe for concurrent use: the site's rules
// collect into the running scrape's batch through AddHeadline.
type DefaultScraper struct {
	Name string
	URL  string

	schedule Schedule
	rules    []headlineRule
	batch    models.Headlines
}

// NewScraper takes the refresh interval in minutes.
func NewScraper(name, url string, minRefreshInterval, maxRefreshInterval int, offHours []int) *DefaultScraper {
	return &DefaultScraper{
		Name: name,
		URL:  url,
		schedule: Schedule{
			MinInterval: time.Duration(minRefreshInterval) * time.Minute,
			MaxInterval: time.Duration(maxRefreshInterval) * time.Minute,
			OffHours:    offHours,
		},
	}
}

func (s *DefaultScraper) String() string {
	return s.Name
}

func (s *DefaultScraper) Schedule() Schedule {
	return s.schedule
}

// ScrapeHeadline registers a rule: callback runs for every element matching selector on the
// front page, and calls AddHeadline for each headline it finds.
func (s *DefaultScraper) ScrapeHeadline(selector string, callback colly.HTMLCallback) {
	s.rules = append(s.rules, headlineRule{selector: selector, callback: callback})
}

// AddHeadline adds a headline to the running scrape, skipping one without a title.
func (s *DefaultScraper) AddHeadline(h models.Headline) {
	title := strings.TrimSpace(h.Title)
	if title == "" {
		return
	}
	h.Title = title
	h.URL = strings.TrimSpace(h.URL)
	s.batch = append(s.batch, h)
}

// Scrape fetches the front page once and returns its headlines in page order, newest first on
// every site scraped today.
func (s *DefaultScraper) Scrape(ctx context.Context) (models.Headlines, error) {
	c := colly.NewCollector(colly.StdlibContext(ctx))
	c.SetRequestTimeout(requestTimeout)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", headlinesUserAgent)
	})
	for _, rule := range s.rules {
		c.OnHTML(rule.selector, rule.callback)
	}

	s.batch = nil
	defer func() { s.batch = nil }()
	if err := c.Visit(s.URL); err != nil {
		return nil, err
	}
	return s.batch, nil
}

// ScrapeStoryFrom fetches a story page and joins the childElement paragraphs inside element into
// sanitized HTML. With html false each paragraph keeps its text and images; with html true it
// keeps its inner HTML, which the sanitizer reduces to plain formatting, images and safe links.
func (s *DefaultScraper) ScrapeStoryFrom(ctx context.Context, url, element, childElement string, html bool) (string, error) {
	var story strings.Builder

	c := colly.NewCollector(colly.StdlibContext(ctx))
	c.UserAgent = storyUserAgent
	c.SetRequestTimeout(requestTimeout)
	c.OnHTML(element, func(e *colly.HTMLElement) {
		e.ForEach(childElement, func(_ int, el *colly.HTMLElement) {
			unwrapNoscript(el.DOM)
			inner, err := el.DOM.Html()
			if err != nil {
				return
			}
			contents := strings.TrimSpace(inner)
			if !html {
				contents = paragraphText(inner)
			}
			if contents != "" {
				story.WriteString("<p>" + contents + "</p>")
			}
		})
	})

	if err := c.Visit(url); err != nil {
		return "", err
	}
	return SanitizeStory(story.String()), nil
}

// unwrapNoscript turns <noscript> content back into elements. Lazy-loading sites put the real
// image there, and the HTML parser keeps that content as raw text.
func unwrapNoscript(s *goquery.Selection) {
	s.Find("noscript").Each(func(_ int, noscript *goquery.Selection) {
		noscript.ReplaceWithHtml(noscript.Text())
	})
}
