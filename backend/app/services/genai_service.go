package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/utils"
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

	// stubSummary is what llm_provider "stub" answers, so tests never reach a paid API.
	stubSummary = "Sažetak vijesti."
)

type GenAIService struct {
	raptor.Service

	provider string
	openai   openai.Client
}

// Setup reads llm_provider: "openai" (the default) needs openai_key and fails boot without it,
// rather than on a user's first summary; "stub" needs nothing.
func (s *GenAIService) Setup() error {
	s.provider = s.Config.AppString("llm_provider", "openai")
	switch s.provider {
	case "stub":
		return nil
	case "openai":
		key := s.Config.AppString("openai_key", "")
		if key == "" {
			return errors.New("openai_key is not set in app config")
		}
		s.openai = openai.NewClient(option.WithAPIKey(key))
		return nil
	default:
		return fmt.Errorf("unknown llm_provider %q (want openai or stub)", s.provider)
	}
}

// Summarize returns a plain-text Croatian summary of story HTML. ctx is the request's, so a
// client that leaves cancels the call; the timeout bounds a hung one.
func (s *GenAIService) Summarize(ctx context.Context, story string) (string, error) {
	text := utils.StoryText(story)
	if s.provider == "stub" {
		return stubSummary, nil
	}

	ctx, cancel := context.WithTimeout(ctx, summarizeTimeout)
	defer cancel()

	result, err := s.openai.Responses.New(ctx, responses.ResponseNewParams{
		Model:        summarizeModel,
		Instructions: openai.String(summarizePrompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(text),
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
	return summary, nil
}
