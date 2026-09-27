package models

import (
	"time"

	"github.com/uptrace/bun"
)

type (
	Headlines = []Headline
	Headline  struct {
		bun.BaseModel `bun:"table:headlines,alias:headlines"`

		ID          int64     `bun:"id,pk,autoincrement" json:"id"`
		SourceID    int64     `bun:"source_id,notnull" json:"sourceId"`
		Title       string    `bun:"title,notnull" json:"title"`
		URL         string    `bun:"url,notnull" json:"url"`
		PublishedAt time.Time `bun:"published_at,nullzero,notnull,default:current_timestamp" json:"publishedAt"`

		CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"createdAt"`
		UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updatedAt"`

		Source *Source `bun:"rel:belongs-to,join:source_id=id" json:"source,omitempty"`
	}
)

// A headline is written by the scrapers and only read by clients, so its response carries
// publishedAt, the one timestamp the UI shows. The source is always loaded, so it is a value.
type (
	HeadlineResponses = []HeadlineResponse
	HeadlineResponse  struct {
		ID          int64          `json:"id"`
		Title       string         `json:"title"`
		URL         string         `json:"url"`
		PublishedAt time.Time      `json:"publishedAt"`
		Source      SourceResponse `json:"source"`
	}
)

func NewHeadlineResponse(h *Headline) HeadlineResponse {
	res := HeadlineResponse{ID: h.ID, Title: h.Title, URL: h.URL, PublishedAt: h.PublishedAt}
	if h.Source != nil {
		res.Source = NewSourceResponse(h.Source)
	}
	return res
}

func NewHeadlineResponses(headlines Headlines) HeadlineResponses {
	res := make(HeadlineResponses, len(headlines))
	for i := range headlines {
		res[i] = NewHeadlineResponse(&headlines[i])
	}
	return res
}

type HeadlineCountResponse struct {
	Count int `json:"count"`
}
