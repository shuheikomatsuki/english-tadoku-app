package service

import (

)

type ILLMService interface {
	GenerateStory(prompt string) (string, error)
}

type LLMService struct {
	APIKey string
}

func NewLLMService(apiKey string) ILLMService {
	return &LLMService{APIKey: apiKey}
}

func (s *LLMService) GenerateStory(prompt string) (string, error) {
	// ここに Gemini API を呼び出すロジックを実装
	content := `This is a generated story based on the prompt: "` + prompt + `".`
	return content, nil
}