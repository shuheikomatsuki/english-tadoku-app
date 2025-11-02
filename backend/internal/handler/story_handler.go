package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"
)

type IStoryHandler interface {
	GenerateStory(e echo.Context) error
	GetStories(e echo.Context) error
	GetStory(e echo.Context) error
	DeleteStory(e echo.Context) error
	UpdateStory(e echo.Context) error
	MarkStoryAsRead(e echo.Context) error
	UndoLastRead(e echo.Context) error
}

type StoryHandler struct {
	// StoryRepo  repository.IStoryRepository
	// LLMService service.ILLMService
	StoryService service.IStoryService
}

type GetStoriesResponse struct {
	Stories     []*model.Story `json:"stories"`
	TotalCount  int            `json:"total_count"`
	TotalPages  int            `json:"total_pages"`
	CurrentPage int            `json:"current_page"`
}

type UpdateStoryRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
}

type StoryDetailResponse struct {
	model.Story
	ReadCount int `json:"read_count"`
}

func NewStoryHandler(storyService service.IStoryService) IStoryHandler {
	return &StoryHandler{
		StoryService: storyService,
	}
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

	// (簡易バリデーション)
	if req.Prompt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "prompt is required"})
	}

	story, err := h.StoryService.GenerateStory(userID, req.Prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate story content"})
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

	paginatedResult, err := h.StoryService.GetStories(userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	res := GetStoriesResponse{
		Stories:     paginatedResult.Stories,
		TotalCount:  paginatedResult.TotalCount,
		TotalPages:  paginatedResult.TotalPages,
		CurrentPage: paginatedResult.CurrentPage,
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

	storyDetail, err := h.StoryService.GetStory(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrStoryNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	res := StoryDetailResponse{
		Story:     storyDetail.Story,
		ReadCount: storyDetail.ReadCount,
	}

	return c.JSON(http.StatusOK, res)
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

	err = h.StoryService.DeleteStory(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrStoryNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		if errors.Is(err, service.ErrForbidden) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "you do not have permission to delete this story"})
		}
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

	updatedStory, err := h.StoryService.UpdateStoryTitle(id, userID, req.Title)
	if err != nil {
		if errors.Is(err, service.ErrStoryNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update story"})
	}

	return c.JSON(http.StatusOK, updatedStory)
}

func (h *StoryHandler) MarkStoryAsRead(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	// JWTからユーザーIDを取得
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	err = h.StoryService.MarkStoryAsRead(id, userID)
	if err != nil {
		if errors.Is(err, service.ErrStoryNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "story not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create reading record"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Story marked as read successfully"})
}

func (h *StoryHandler) UndoLastRead(c echo.Context) error {
	storyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid story id"})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	err = h.StoryService.UndoLastRead(storyID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNoReadingRecord) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No reading record found to undo."})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete reading record"})
	}

	return c.JSON(http.StatusNoContent, nil)
}