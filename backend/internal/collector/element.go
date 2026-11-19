package collector

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLElement is one element matched by an OnHTML selector or a ForEach.
type HTMLElement struct {
	// Name is the tag name, such as "a" or "div".
	Name string
	// Text is the element's text content, untrimmed.
	Text string
	// DOM is the element as a goquery selection, for anything the helpers don't cover.
	DOM *goquery.Selection
}

// newHTMLElement takes a single-node selection, as goquery's Each yields.
func newHTMLElement(s *goquery.Selection) *HTMLElement {
	return &HTMLElement{
		Name: goquery.NodeName(s),
		Text: s.Text(),
		DOM:  s,
	}
}

// Attr returns the element's attribute k as written, or "" when it has none.
func (e *HTMLElement) Attr(k string) string {
	v, _ := e.DOM.Attr(k)
	return v
}

// ChildText returns the trimmed, concatenated text of every descendant matching selector.
func (e *HTMLElement) ChildText(selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).Text())
}

// ChildAttr returns the trimmed attribute attr of the first descendant matching selector, or ""
// when it has none.
func (e *HTMLElement) ChildAttr(selector, attr string) string {
	v, _ := e.DOM.Find(selector).Attr(attr)
	return strings.TrimSpace(v)
}

// ForEach runs fn for every descendant matching selector, in document order, with its index.
func (e *HTMLElement) ForEach(selector string, fn func(int, *HTMLElement)) {
	e.DOM.Find(selector).Each(func(i int, s *goquery.Selection) {
		fn(i, newHTMLElement(s))
	})
}
