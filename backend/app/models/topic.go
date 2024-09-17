package models

type Topic struct {
	ID      int64     `json:"id" bun:",pk,autoincrement"`
	Name    string    `json:"name"`
	Sources []*Source `json:"-" bun:"rel:has-many,join:id=topic_id"`
}

type Topics []Topic
