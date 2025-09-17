package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func getUserIDFromContext(c echo.Context) (int, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("failed to get user from context")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("failed to get claims from token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in claims or is not a number")
	}

	return int(userIDFloat), nil

}

func (h *StoryHandler) GenerateStory(c echo.Context) error {
	// TODO: JWTミドルウェアからユーザーIDを取得する
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	// userID := 123

	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	content, err := h.LLMService.GenerateStory(req.Prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate story content"})
	}

	story := &model.Story{
		UserID: userID,
		Title: req.Prompt,
		Content: content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.StoryRepo.CreateStory(story); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save story"})
	}

	return c.JSON(http.StatusCreated, story)

	// return nil
}

func (h *StoryHandler) GetStories(c echo.Context) error {
	// TODO: JWTミドルウェアからユーザーIDを取得する
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	// userID := 123

	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offsetStr := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	stories, err := h.StoryRepo.GetUserStories(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusOK, stories)

	// return nil
}

func (h *StoryHandler) GetStory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	// TODO: JWTミドルウェアからユーザーIDを取得する
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	// userID := 123 // テスト用の仮ID

	story, err := h.StoryRepo.GetUserStory(id, userID)
	if err != nil {
		if err == sql.ErrNoRows || err == http.ErrMissingFile {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusOK, story)

	// return nil
}

func (h *StoryHandler) DeleteStory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	// TODO: JWTミドルウェアからユーザーIDを取得する
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	
	// TODO: 削除する story が本当にこのユーザーのものかを確認する。
	story, err := h.StoryRepo.GetUserStory(id, userID)
	if err != nil {
		if err == sql.ErrNoRows || err == http.ErrMissingFile {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	if story.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "you do not have permission to delete this story"})
	}

	if err := h.StoryRepo.DeleteStory(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete story"})
	}

	return c.JSON(http.StatusNoContent, nil)

	// return nil
}