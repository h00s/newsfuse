package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Story struct {
	bun.BaseModel `bun:"table:stories,alias:stories"`

	ID         int64 `bun:"id,pk,autoincrement" json:"id"`
	HeadlineID int64 `bun:"headline_id,notnull" json:"headlineId"`
	// Content is sanitized HTML; Summary is plain text, empty until the story is summarized.
	Content string `bun:"content,notnull" json:"content"`
	Summary string `bun:"summary,notnull" json:"summary"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

type StoryResponse struct {
	ID         int64  `json:"id"`
	HeadlineID int64  `json:"headlineId"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
}

func NewStoryResponse(s *Story) StoryResponse {
	return StoryResponse{ID: s.ID, HeadlineID: s.HeadlineID, Content: s.Content, Summary: s.Summary}
}
