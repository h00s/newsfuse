package models

type Source struct {
	ID          int64       `json:"id" bun:",pk,autoincrement"`
	Name        string      `json:"name"`
	IsScrapable bool        `json:"isScrapable"`
	TopicID     int64       `json:"topicId" bun:",notnull"`
	Headlines   []*Headline `json:"-" bun:"rel:has-many,join:id=source_id"`
}

type Sources []Source
