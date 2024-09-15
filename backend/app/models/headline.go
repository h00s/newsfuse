package models

import (
	"time"
)

type Headline struct {
	ID          int64     `json:"id" bun:",pk,autoincrement"`
	Title       string    `json:"title"`
	URL         string    `json:"url" bun:",unique"`
	SourceID    uint      `json:"source_id"`
	Source      *Source   `json:"-" bun:"rel:belongs-to,join:source_id=id"`
	Story       *Story    `json:"-"`
	PublishedAt time.Time `json:"published_at" bun:",nullzero,notnull,default:current_timestamp"`
	CreatedAt   time.Time `json:"-" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `json:"-" bun:",nullzero,notnull,default:current_timestamp"`
}

type Headlines []Headline
