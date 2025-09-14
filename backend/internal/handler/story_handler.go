// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"
)

type IStoryHandler interface {
	GenerateStory(e echo.Context) error
	GetStories(e echo.Context) error
	GetStory(e echo.Context) error
	DeleteStory(e echo.Context) error
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

func (h *StoryHandler) GenerateStory(e echo.Context) error {
	return nil
}

func (h *StoryHandler) GetStories(e echo.Context) error {
	return nil
}

func (h *StoryHandler) GetStory(e echo.Context) error {
	return nil
}

func (h *StoryHandler) DeleteStory(e echo.Context) error {
	return nil
}