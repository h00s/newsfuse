package models

import (
	"time"

	"gorm.io/gorm"
)

type Headline struct {
	gorm.Model
	Title       string
	URL         string    `gorm:"index"`
	PublishedAt time.Time `gorm:"index"`
	Source      string
}

type Headlines []Headline
