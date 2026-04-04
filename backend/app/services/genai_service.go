package services

import (
	"context"
	"strings"

	"github.com/go-raptor/raptor/v4"
	"google.golang.org/genai"
)

type GenAIService struct {
	raptor.Service

	genai *genai.Client
}

func (s *GenAIService) Setup() error {
	var err error
	s.genai, err = genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  s.Config.AppConfig["gemini_key"],
		Backend: genai.BackendGeminiAPI,
	})
	return err
}

func (s *GenAIService) Summarize(story string) (string, error) {
	replacer := strings.NewReplacer("<p>", "", "</p>", "")
	story = replacer.Replace(story)
	/* if len(story) > 2500 {
		story = story[:2500]
	} */
	content := "Napravi sažetak vijesti bez navoda da se radi o sažetku. Odgovor mora sadržavati samo tekst sažetka vijesti i to do 600 znakova na hrvatskom jeziku: " + story

	result, err := s.genai.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash-lite",
		genai.Text(content),
		nil,
	)

	if err == nil {
		return "<p>" + result.Text() + "</p>", nil
	}

	return "", err
}
