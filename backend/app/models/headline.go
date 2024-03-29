package models

import (
	"time"
)

type Headline struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	URL         string    `gorm:"unique" json:"url"`
	SourceID    uint      `json:"source_id"`
	Source      Source    `json:"-"`
	Story       Story     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	PublishedAt time.Time `gorm:"index" json:"published_at"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Headlines []Headline
