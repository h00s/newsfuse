package models

import "time"

type Headline struct {
	Title       string
	Source      Source
	URL         string
	PublishedAt time.Time
}

type Headlines []Headline
