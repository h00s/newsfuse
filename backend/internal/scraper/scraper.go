// Package scraper is the base every news site builds on: a front page yields headlines, and a
// headline's page yields its story.
package scraper

import (
	"context"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/h00s/newsfuse/internal/collector"
	"github.com/h00s/newsfuse/internal/sanitize"
)

const (
	requestTimeout = 15 * time.Second
	// userAgent is a regular browser's, for headlines and stories alike. Sites behind bot
	// protection reject a Googlebot user agent that doesn't come from Google's network: Hacker
	// News answers 419, for one.
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:151.0) Gecko/20100101 Firefox/151.0"
)

// Headline is one headline on a site's front page. The app stores it under the site's source.
type Headline struct {
	Title string
	URL   string
}

// Scraper is one news site: its front page yields headlines, and each headline's page yields its
// story. String names the site in logs.
type Scraper interface {
	String() string
	Schedule() Schedule
	Scrape(ctx context.Context) ([]Headline, error)
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
	callback collector.HTMLCallback
}

// Base does the fetching; each site embeds it, adds its headline rules with ScrapeHeadline
// and picks its story selectors in ScrapeStory. Scrape is not safe for concurrent use: the site's rules
// collect into the running scrape's batch through AddHeadline.
type Base struct {
	Name string
	URL  string

	schedule Schedule
	rules    []headlineRule
	batch    []Headline
}

// New takes the refresh interval in minutes.
func New(name, url string, minRefreshInterval, maxRefreshInterval int, offHours []int) *Base {
	return &Base{
		Name: name,
		URL:  url,
		schedule: Schedule{
			MinInterval: time.Duration(minRefreshInterval) * time.Minute,
			MaxInterval: time.Duration(maxRefreshInterval) * time.Minute,
			OffHours:    offHours,
		},
	}
}

func (s *Base) String() string {
	return s.Name
}

func (s *Base) Schedule() Schedule {
	return s.schedule
}

// ScrapeHeadline registers a rule: callback runs for every element matching selector on the
// front page, and calls AddHeadline for each headline it finds.
func (s *Base) ScrapeHeadline(selector string, callback collector.HTMLCallback) {
	s.rules = append(s.rules, headlineRule{selector: selector, callback: callback})
}

// AddHeadline adds a headline to the running scrape, skipping one without a title.
func (s *Base) AddHeadline(h Headline) {
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
func (s *Base) Scrape(ctx context.Context) ([]Headline, error) {
	c := collector.New(collector.UserAgent(userAgent), collector.Timeout(requestTimeout))
	for _, rule := range s.rules {
		c.OnHTML(rule.selector, rule.callback)
	}

	s.batch = nil
	defer func() { s.batch = nil }()
	if err := c.Visit(ctx, s.URL); err != nil {
		return nil, err
	}
	return s.batch, nil
}

// ScrapeStoryFrom fetches a story page and joins the childElement paragraphs inside element into
// sanitized HTML. With html false each paragraph keeps its text and images; with html true it
// keeps its inner HTML, which the sanitizer reduces to plain formatting, images and safe links.
func (s *Base) ScrapeStoryFrom(ctx context.Context, url, element, childElement string, html bool) (string, error) {
	var story strings.Builder

	c := collector.New(collector.UserAgent(userAgent), collector.Timeout(requestTimeout))
	c.OnHTML(element, func(e *collector.HTMLElement) {
		e.ForEach(childElement, func(_ int, el *collector.HTMLElement) {
			unwrapNoscript(el.DOM)
			inner, err := el.DOM.Html()
			if err != nil {
				return
			}
			contents := strings.TrimSpace(inner)
			if !html {
				contents = sanitize.ParagraphText(inner)
			}
			if contents != "" {
				story.WriteString("<p>" + contents + "</p>")
			}
		})
	})

	if err := c.Visit(ctx, url); err != nil {
		return "", err
	}
	return sanitize.Story(story.String()), nil
}

// unwrapNoscript turns <noscript> content back into elements. Lazy-loading sites put the real
// image there, and the HTML parser keeps that content as raw text.
func unwrapNoscript(s *goquery.Selection) {
	s.Find("noscript").Each(func(_ int, noscript *goquery.Selection) {
		noscript.ReplaceWithHtml(noscript.Text())
	})
}
