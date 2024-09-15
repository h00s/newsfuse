package models

type Source struct {
	ID          int64       `json:"id" bun:",pk,autoincrement"`
	Name        string      `json:"name"`
	IsScrapable bool        `json:"is_scrapable"`
	TopicID     int64       `json:"topic_id" bun:",notnull"`
	Headlines   []*Headline `json:"-" bun:"rel:has-many,join:id=source_id"`
}

type Sources []Source
