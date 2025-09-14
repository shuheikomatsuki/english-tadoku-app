// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"
)

type IStoryHandler interface {
	// CreateStory() error
	// GetStories() error
	// GetStory() error
	// DeleteStory() error
}

type StoryHandler struct {
	StoryRepo repository.IStoryRepository
	LLMService service.ILLMService
}

func NewStoryHandler(storyRepo repository.IStoryRepository, llmService service.ILLMService) IStoryHandler {
	return &StoryHandler{
		StoryRepo: storyRepo,
		LLMService: llmService,
	}
}