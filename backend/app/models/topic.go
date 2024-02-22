package models

import (
	"time"
)

type Topic struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Sources   []Source  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

type Topics []Topic
