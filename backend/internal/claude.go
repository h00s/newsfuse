package internal

import (
	"strings"

	"github.com/madebywelch/anthropic-go/v2/pkg/anthropic"
)

type Claude struct {
	client *anthropic.Client
}

func NewClaude(apiKey string) (*Claude, error) {
	if client, err := anthropic.NewClient(apiKey); err == nil {
		return &Claude{
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (c *Claude) Summarize(story string) (string, error) {
	replacer := strings.NewReplacer("<p>", "", "</p>", "")
	story = replacer.Replace(story)
	if len(story) > 2500 {
		story = story[:2500]
	}
	content := "Napravi sažetak vijesti. Odgovori samo sa sažetkom vijesti do 600 znakova bez navoda da se radi o sažetku: " + story
	request := anthropic.NewMessageRequest(
		[]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock(content)}}},
		anthropic.WithModel[anthropic.MessageRequest](anthropic.Claude3Haiku),
		anthropic.WithMaxTokens[anthropic.MessageRequest](1000),
	)

	response, err := c.client.Message(request)
	if err != nil {
		return "", err
	}

	return "<p>" + response.Content[0].Text + "</p>", nil
}
