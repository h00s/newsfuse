package internal

import (
	"context"

	gogpt "github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	client *gogpt.Client
}

func NewChatGPT(authToken string) *ChatGPT {
	return &ChatGPT{
		client: gogpt.NewClient(authToken),
	}
}

func (c *ChatGPT) Summarize(story string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		gogpt.ChatCompletionRequest{
			Model: gogpt.GPT3Dot5Turbo0125,
			Messages: []gogpt.ChatCompletionMessage{
				{
					Role:    gogpt.ChatMessageRoleUser,
					Content: "Napravi mi sažetak do 600 znakova sljedeće vijesti: " + story,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
