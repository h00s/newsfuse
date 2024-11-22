package internal

import (
	"context"
	"strings"

	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic"
	"github.com/madebywelch/anthropic-go/v3/pkg/anthropic/client/native"
)

type Claude struct {
	client *native.Client
}

func NewClaude(apiKey string) (*Claude, error) {
	if client, err := native.MakeClient(native.Config{APIKey: apiKey}); err == nil {
		return &Claude{
			client: client,
		}, nil
	} else {
		return nil, err
	}
}

func (c *Claude) Summarize(story string) (string, error) {
	ctx := context.Background()
	replacer := strings.NewReplacer("<p>", "", "</p>", "")
	story = replacer.Replace(story)
	if len(story) > 2500 {
		story = story[:2500]
	}
	content := "Napravi sažetak vijesti bez navoda da se radi o sažetku. Odgovor mora sadržavati samo tekst sažetka vijesti i to do 600 znakova: " + story
	request := anthropic.NewMessageRequest(
		[]anthropic.MessagePartRequest{{Role: "user", Content: []anthropic.ContentBlock{anthropic.NewTextContentBlock(content)}}},
		anthropic.WithModel[anthropic.MessageRequest](anthropic.Claude35Haiku),
		anthropic.WithMaxTokens[anthropic.MessageRequest](1000),
	)

	response, err := c.client.Message(ctx, request)
	if err != nil {
		return "", err
	}

	return "<p>" + response.Content[0].Text + "</p>", nil
}
