package models

import "time"

type Headline struct {
	ID         int64     `json:"id"`
	SourceID   int       `json:"source_id"`
	Title      string    `json:"title"`
	URL        string    `json:"link"`
	InsertedAt time.Time `json:"inserted_at"`
}
