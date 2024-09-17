package models

import "time"

type Story struct {
	ID         int64     `json:"id" bun:",pk,autoincrement"`
	HeadlineID int64     `json:"headline_id" bun:",unique,notnull"`
	Headline   *Headline `json:"-" bun:"rel:belongs-to,join:headline_id=id"` // cascade
	Summary    string    `json:"summary" bun:",type:text"`
	Content    string    `json:"content" bun:",type:text"`
	CreatedAt  time.Time `json:"-" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt  time.Time `json:"-" bun:",nullzero,notnull,default:current_timestamp"`
}
