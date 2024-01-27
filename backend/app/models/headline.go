package models

import "time"

type Headline struct {
	Title string
	//Source      Source
	URL         string
	PublishedAt time.Time
	Source      string
}

type Headlines []Headline
