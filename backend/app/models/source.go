package models

type Source struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	IsScrapable bool       `json:"is_scrapable"`
	TopicID     uint       `gorm:"not null" json:"topic_id"`
	Headlines   []Headline `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

type Sources []Source
