package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

type ILLMService interface {
	GenerateStory(prompt string) (string, error)
}

type LLMService struct {
	APIKey string
	client *genai.Client
}

func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func NewLLMService(apiKey string) (ILLMService, error) {
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &LLMService{
		APIKey: apiKey,
		client: client,
	}, nil
}

func (s *LLMService) GenerateStory(prompt string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("genai client is not initialized")
	}

instructionalPrompt := fmt.Sprintf(
`Write a clear, factual explanation in English based on the user's prompt.
The user's prompt may be written in Japanese or English.
Always write the output in English.
Use standard Markdown for paragraphs and lists where appropriate.
Do not write a story, narrative, or fictional content.
Return only the Markdown content, without explanations or notes outside the text.

--- USER PROMPT START ---
%s
--- USER PROMPT END ---`,
prompt,
)

	// API 呼び出しにタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	const model = "gemini-2.5-flash-lite"

	result, err := s.client.Models.GenerateContent(ctx, model, genai.Text(instructionalPrompt), nil)
	if err != nil {
		// エラーをログに出す（500 の原因調査用）
		fmt.Printf("Gemini API error: %v\n", err)
		return "", fmt.Errorf("generate content failed: %w", err)
	}

	if result == nil {
		return "", fmt.Errorf("gemini returned nil response")
	}

	text := result.Text()
	if text == "" {
		return "", fmt.Errorf("gemini returned empty text response: %+v", result)
	}

	return text, nil
}
