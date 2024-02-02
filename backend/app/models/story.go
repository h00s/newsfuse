package models

import "github.com/go-raptor/raptor"

type Story struct {
	raptor.Model
	Summary    string `gorm:"type:text"`
	HeadlineID uint
}
