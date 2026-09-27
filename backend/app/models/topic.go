// Package models holds the Bun models and the API's response DTOs.
package models

import "github.com/uptrace/bun"

type (
	Topics = []Topic
	Topic  struct {
		bun.BaseModel `bun:"table:topics,alias:topics"`

		ID   int64  `bun:"id,pk,autoincrement" json:"id"`
		Name string `bun:"name,notnull" json:"name"`
	}
)

// Topics are seeded reference data: no timestamps, nothing consumes them.
type (
	TopicResponses = []TopicResponse
	TopicResponse  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func NewTopicResponse(t *Topic) TopicResponse {
	return TopicResponse{ID: t.ID, Name: t.Name}
}

func NewTopicResponses(topics Topics) TopicResponses {
	res := make(TopicResponses, len(topics))
	for i := range topics {
		res[i] = NewTopicResponse(&topics[i])
	}
	return res
}
