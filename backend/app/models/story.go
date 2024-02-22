package models

import "time"

type Story struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	HeadlineID uint      `gorm:"unique;not null" json:"headline_id"`
	Summary    string    `gorm:"type:text" json:"summary"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
