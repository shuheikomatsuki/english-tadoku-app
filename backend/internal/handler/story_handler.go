// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
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
	// TODO: JWTミドルウェアからユーザーIDを取得する
	// claims := e.Get("user").(*jwt.Token).Claims.(*tokenClaims)
	// userID := claims.UserID
	userID := 123

	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	content, err := h.LLMService.GenerateStory(req.Prompt)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate story content"})
	}

	story := &model.Story{
		UserID: userID,
		Title: req.Prompt,
		Content: content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.StoryRepo.CreateStory(story); err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save story"})
	}

	return e.JSON(http.StatusCreated, story)

	// return nil
}

func (h *StoryHandler) GetStories(e echo.Context) error {
	return nil
}

func (h *StoryHandler) GetStory(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	// TODO: JWTミドルウェアからユーザーIDを取得する
	userID := 123 // テスト用の仮ID

	story, err := h.StoryRepo.GetUserStory(id, userID)
	if err != nil {
		if err == sql.ErrNoRows || err == http.ErrMissingFile {
			return e.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return e.JSON(http.StatusOK, story)

	// return nil
}

func (h *StoryHandler) DeleteStory(e echo.Context) error {
	return nil
}