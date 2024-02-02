package models

import (
	"time"

	"github.com/go-raptor/raptor"
)

type Headline struct {
	raptor.Model
	Title       string
	URL         string    `gorm:"index"`
	PublishedAt time.Time `gorm:"index"`
	SourceID    uint
	Story       Story `json:"-"`
}

type Headlines []Headline
