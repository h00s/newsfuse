package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

const (
	summarizeModel   = openai.ChatModelGPT5_6Luna
	summarizeTimeout = 60 * time.Second
	summarizePrompt  = "Napravi sažetak vijesti bez navoda da se radi o sažetku. " +
		"Odgovor mora sadržavati samo tekst sažetka vijesti i to do 600 znakova na hrvatskom jeziku."
)

type GenAIService struct {
	raptor.Service

	openai openai.Client
}

func (s *GenAIService) Setup() error {
	key := s.Config.AppConfig["openai_key"]
	if key == "" {
		return errors.New("openai_key is not set in app config")
	}
	s.openai = openai.NewClient(option.WithAPIKey(key))
	return nil
}

func (s *GenAIService) Summarize(story string) (string, error) {
	replacer := strings.NewReplacer("<p>", "", "</p>", "")
	story = replacer.Replace(story)

	ctx, cancel := context.WithTimeout(context.Background(), summarizeTimeout)
	defer cancel()

	result, err := s.openai.Responses.New(ctx, responses.ResponseNewParams{
		Model:        summarizeModel,
		Instructions: openai.String(summarizePrompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(story),
		},
		Reasoning: shared.ReasoningParam{
			Effort: shared.ReasoningEffortLow,
		},
		Text: responses.ResponseTextConfigParam{
			Verbosity: responses.ResponseTextConfigVerbosityLow,
		},
	})
	if err != nil {
		return "", err
	}

	summary := strings.TrimSpace(result.OutputText())
	if summary == "" {
		return "", errors.New("empty summary returned")
	}

	return "<p>" + summary + "</p>", nil
}
