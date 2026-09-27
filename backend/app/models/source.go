package models

import "github.com/uptrace/bun"

type (
	Sources = []Source
	Source  struct {
		bun.BaseModel `bun:"table:sources,alias:sources"`

		ID          int64  `bun:"id,pk,autoincrement" json:"id"`
		TopicID     int64  `bun:"topic_id,notnull" json:"topicId"`
		Name        string `bun:"name,notnull" json:"name"`
		IsScrapable bool   `bun:"is_scrapable,notnull" json:"isScrapable"`
	}
)

// Sources are seeded reference data, like topics.
type (
	SourceResponses = []SourceResponse
	SourceResponse  struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		TopicID     int64  `json:"topicId"`
		IsScrapable bool   `json:"isScrapable"`
	}
)

func NewSourceResponse(s *Source) SourceResponse {
	return SourceResponse{ID: s.ID, Name: s.Name, TopicID: s.TopicID, IsScrapable: s.IsScrapable}
}

func NewSourceResponses(sources Sources) SourceResponses {
	res := make(SourceResponses, len(sources))
	for i := range sources {
		res[i] = NewSourceResponse(&sources[i])
	}
	return res
}
