package models

import "github.com/go-raptor/raptor"

type Source struct {
	raptor.Model `json:"-"`
	Name         string
	TopicID      uint       `gorm:"not null"`
	Headlines    []Headline `json:"-"`
}

type Sources []Source
