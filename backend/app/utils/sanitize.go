package utils

import (
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
)

// allowImages permits img with an absolute http(s) src, its alt text and its intrinsic size, which
// lets the browser reserve the space before the image loads.
func allowImages(p *bluemonday.Policy) {
	p.AllowAttrs("src").OnElements("img")
	p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("img")
	p.AllowAttrs("width", "height").Matching(bluemonday.Integer).OnElements("img")
}

// storyPolicy keeps what an article needs to read well and nothing that can run: paragraphs,
// emphasis, lists, quotes, images and absolute http(s) links, which open in a new tab without a
// referrer.
var storyPolicy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "br", "b", "strong", "i", "em", "ul", "ol", "li", "blockquote")
	p.AllowAttrs("href").OnElements("a")
	allowImages(p)
	p.AllowURLSchemes("http", "https")
	p.RequireParseableURLs(true)
	p.RequireNoFollowOnLinks(true)
	p.RequireNoReferrerOnLinks(true)
	p.AddTargetBlankToFullyQualifiedLinks(true)
	return p
}()

// paragraphTextPolicy is a paragraph's text and images, without links or formatting: what a
// scraper in text mode keeps.
var paragraphTextPolicy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	allowImages(p)
	p.AllowURLSchemes("http", "https")
	p.RequireParseableURLs(true)
	return p
}()

var textPolicy = bluemonday.StrictPolicy()

// SanitizeStory makes scraped story HTML safe to render as HTML. An image whose source was
// removed (a lazy-loading placeholder, an unsafe URL) goes too, and so does a paragraph left empty.
func SanitizeStory(content string) string {
	clean := storyPolicy.Sanitize(content)
	if !strings.Contains(clean, "<img") {
		return clean
	}
	return dropEmptyImages(clean)
}

func dropEmptyImages(fragment string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fragment))
	if err != nil {
		return fragment // already sanitized: at worst an image without a source stays
	}
	doc.Find("img:not([src])").Remove()
	doc.Find("p").FilterFunction(func(_ int, p *goquery.Selection) bool {
		return strings.TrimSpace(p.Text()) == "" && p.Find("img").Length() == 0
	}).Remove()
	out, err := doc.Find("body").Html()
	if err != nil {
		return fragment
	}
	return out
}

// paragraphText is a text-mode paragraph: its text, escaped, and its images.
func paragraphText(innerHTML string) string {
	return strings.TrimSpace(paragraphTextPolicy.Sanitize(innerHTML))
}

// StoryText reduces story HTML to plain text, one paragraph per line, for the summarizer.
func StoryText(content string) string {
	content = strings.ReplaceAll(content, "</p>", "</p>\n")
	return strings.TrimSpace(html.UnescapeString(textPolicy.Sanitize(content)))
}
