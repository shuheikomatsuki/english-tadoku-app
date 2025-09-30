package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	// "time"

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
	UpdateStory(e echo.Context) error
}

type StoryHandler struct {
	StoryRepo repository.IStoryRepository
	LLMService service.ILLMService
}

type GetStoriesResponse struct {
	Stories     []*model.Story  `json:"stories"`
	TotalCount  int             `json:"total_count"`
	TotalPages  int             `json:"total_pages"`
	CurrentPage int             `json:"current_page"`
}

type UpdateStoryRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
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

	claims, ok := user.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, fmt.Errorf("failed to get claims from token")
	}

	return claims.UserID, nil
}

func (h *StoryHandler) GenerateStory(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

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

	wordCount := len(strings.Fields(content))

	story := &model.Story{
		UserID: userID,
		Title: req.Prompt,
		Content: content,
		WordCount: wordCount,
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}

	if err := h.StoryRepo.CreateStory(story); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save story"})
	}

	return c.JSON(http.StatusCreated, story)
}

func (h *StoryHandler) GetStories(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	// --- クエリパラメータの解釈 ---
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	totalCount, err := h.StoryRepo.CountUserStories(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	stories, err := h.StoryRepo.GetUserStories(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	if stories == nil {
		stories = []*model.Story{}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	res := GetStoriesResponse{
		Stories:    stories,
		TotalCount: totalCount,
		TotalPages: totalPages,
		CurrentPage: page,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *StoryHandler) GetStory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	story, err := h.StoryRepo.GetUserStory(id, userID)
	if err != nil {
		if err == sql.ErrNoRows || err == http.ErrMissingFile {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	return c.JSON(http.StatusOK, story)
}

func (h *StoryHandler) DeleteStory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	
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
}

func (h *StoryHandler) UpdateStory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}
	
	var req UpdateStoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	updatedStory, err := h.StoryRepo.UpdateStoryTitle(id, userID, req.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update story"})
	}

	return c.JSON(http.StatusOK, updatedStory)
}