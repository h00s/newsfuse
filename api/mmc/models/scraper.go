package models

type Scraper interface {
	GetHeadlines() ([]Headline, error)
	GetInterval() int
}
