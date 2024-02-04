package models

import "github.com/go-raptor/raptor"

type Story struct {
	raptor.Model
	Summary    string `gorm:"type:text"`
	Content    string `gorm:"type:text"`
	HeadlineID uint   `gorm:"unique;not null"`
}
