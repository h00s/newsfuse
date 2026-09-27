package utils

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// storyPolicy keeps what an article needs to read well and nothing that can run: paragraphs,
// emphasis, lists, quotes and absolute http(s) links, which open in a new tab without a referrer.
var storyPolicy = func() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "br", "b", "strong", "i", "em", "ul", "ol", "li", "blockquote")
	p.AllowAttrs("href").OnElements("a")
	p.AllowURLSchemes("http", "https")
	p.RequireParseableURLs(true)
	p.RequireNoFollowOnLinks(true)
	p.RequireNoReferrerOnLinks(true)
	p.AddTargetBlankToFullyQualifiedLinks(true)
	return p
}()

var textPolicy = bluemonday.StrictPolicy()

// SanitizeStory makes scraped story HTML safe to render as HTML.
func SanitizeStory(content string) string {
	return storyPolicy.Sanitize(content)
}

// StoryText reduces story HTML to plain text, one paragraph per line, for the summarizer.
func StoryText(content string) string {
	content = strings.ReplaceAll(content, "</p>", "</p>\n")
	return strings.TrimSpace(html.UnescapeString(textPolicy.Sanitize(content)))
}
