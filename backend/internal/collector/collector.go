// Package collector is a small HTML scraping client modeled on colly: register callbacks for CSS
// selectors with OnHTML, then Visit a page to fetch it and run them over the parsed document.
// It keeps only what newsfuse's scrapers need, on net/http and goquery.
package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const (
	DefaultTimeout     = 10 * time.Second
	DefaultMaxBodySize = 10 << 20
	DefaultUserAgent   = "newsfuse"
)

// HTMLCallback runs for every element that matches its OnHTML selector.
type HTMLCallback func(*HTMLElement)

// Option configures a Collector in New.
type Option func(*Collector)

// UserAgent sets the User-Agent header sent with every request.
func UserAgent(ua string) Option {
	return func(c *Collector) { c.userAgent = ua }
}

// Timeout limits a whole request, from connecting to reading the last byte of the body.
func Timeout(d time.Duration) Option {
	return func(c *Collector) { c.client.Timeout = d }
}

// MaxBodySize caps how many bytes of a response are read; the rest is dropped, as in colly.
// Zero or less means no cap.
func MaxBodySize(n int64) Option {
	return func(c *Collector) { c.maxBodySize = n }
}

// StatusError is returned by Visit for a response outside 2xx.
type StatusError struct {
	URL        string
	StatusCode int
}

func (e *StatusError) Error() string {
	if text := http.StatusText(e.StatusCode); text != "" {
		return fmt.Sprintf("GET %s: %d %s", e.URL, e.StatusCode, text)
	}
	return fmt.Sprintf("GET %s: status %d", e.URL, e.StatusCode)
}

type htmlCallback struct {
	selector string
	fn       HTMLCallback
}

// Collector fetches pages and runs the registered callbacks over them. Register callbacks before
// visiting; a Collector is not safe for concurrent registration and visiting. Cookies set by a
// site are kept for the Collector's lifetime, so they survive redirects.
type Collector struct {
	userAgent   string
	maxBodySize int64
	client      *http.Client
	callbacks   []htmlCallback
}

// New returns a Collector with colly's defaults: a 10 second timeout, a 10 MB body cap and a
// cookie jar of its own. Connections are pooled across Collectors through http.DefaultTransport.
func New(options ...Option) *Collector {
	jar, _ := cookiejar.New(nil) // never fails without options
	c := &Collector{
		userAgent:   DefaultUserAgent,
		maxBodySize: DefaultMaxBodySize,
		client:      &http.Client{Jar: jar, Timeout: DefaultTimeout},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// OnHTML registers fn to run for every element matching the goquery selector. Callbacks run in
// registration order, each over all of its matches in document order.
func (c *Collector) OnHTML(selector string, fn HTMLCallback) {
	c.callbacks = append(c.callbacks, htmlCallback{selector: selector, fn: fn})
}

// Visit fetches url and runs the HTML callbacks over it. A response outside 2xx is a
// *StatusError; a response that isn't HTML succeeds without running any callback. The body is
// decoded to UTF-8 from the charset its Content-Type or <meta> declares.
func (c *Collector) Visit(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &StatusError{URL: url, StatusCode: resp.StatusCode}
	}

	var body io.Reader = resp.Body
	if c.maxBodySize > 0 {
		body = io.LimitReader(body, c.maxBodySize)
	}
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(c.callbacks) == 0 {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(content)
	}
	if !isHTML(contentType) {
		return nil
	}

	decoded, err := charset.NewReader(bytes.NewReader(content), contentType)
	if err != nil {
		return fmt.Errorf("GET %s: decoding body: %w", url, err)
	}
	doc, err := goquery.NewDocumentFromReader(decoded)
	if err != nil {
		return fmt.Errorf("GET %s: parsing HTML: %w", url, err)
	}

	for _, cb := range c.callbacks {
		doc.Find(cb.selector).Each(func(_ int, s *goquery.Selection) {
			cb.fn(newHTMLElement(s))
		})
	}
	return nil
}

func isHTML(contentType string) bool {
	mediaType, _, _ := strings.Cut(contentType, ";")
	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "text/html", "application/xhtml+xml":
		return true
	}
	return false
}
