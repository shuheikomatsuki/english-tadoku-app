package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/readoku/backend/internal/model"
	"github.com/shuheikomatsuki/readoku/backend/internal/service"
	"github.com/shuheikomatsuki/readoku/backend/internal/timeutil"
)

const (
	testUserID  = 1
	testStoryID = 10
)

func setupTestHandler(t *testing.T) (*MockStoryService, *echo.Echo, *jwt.Token) {
	mockStoryService := new(MockStoryService)
	e := echo.New()
	e.Validator = NewValidator()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(timeutil.NowTokyo().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return mockStoryService, e, token
}

func TestStoryHandler_GetStory(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should return a story", func(t *testing.T) {
		expectedDetail := &service.StoryDetail{
			Story:     *testStory,
			ReadCount: 5,
		}

		mockStoryService.On("GetStory", testStoryID, testUserID).Return(expectedDetail, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		require.NoError(t, h.GetStory(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedResponse StoryDetailResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &receivedResponse))
		assert.Equal(t, expectedDetail.Title, receivedResponse.Title)
		assert.Equal(t, expectedDetail.ReadCount, receivedResponse.ReadCount)

		mockStoryService.AssertExpectations(t)
	})

	t.Run("fail: should return 404 Not Found for non-existent story", func(t *testing.T) {
		nonExistingStoryID := 99
		mockStoryService.On("GetStory", nonExistingStoryID, testUserID).Return(nil, service.ErrStoryNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(nonExistingStoryID))

		require.NoError(t, h.GetStory(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockStoryService.AssertExpectations(t)
	})
}

func TestStoryHandler_GetStories(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should return a list of story", func(t *testing.T) {
		page, limit := 1, 10

		expectedResult := &service.PaginatedStories{
			Stories:     []*model.Story{testStory},
			TotalCount:  1,
			TotalPages:  1,
			CurrentPage: page,
		}

		mockStoryService.On("GetStories", testUserID, page, limit).Return(expectedResult, nil).Once()

		reqURL := fmt.Sprintf("/stories?page=%d&limit=%d", page, limit)
		req := httptest.NewRequest(http.MethodGet, reqURL, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)

		require.NoError(t, h.GetStories(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var response GetStoriesResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Stories, 1)
		assert.Equal(t, expectedResult.TotalCount, response.TotalCount)
		assert.Equal(t, expectedResult.TotalPages, response.TotalPages)
		assert.Equal(t, expectedResult.CurrentPage, response.CurrentPage)

		mockStoryService.AssertExpectations(t)
	})
}

func TestStoryHandler_GenerateStory(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should generate and save a story", func(t *testing.T) {
		prompt := "A story abount Go"
		requestBody := fmt.Sprintf(`{"prompt": "%s"}`, prompt)

		generatedStory := *testStory
		generatedStory.Title = prompt

		mockStoryService.On("GenerateStory", testUserID, prompt).Return(&generatedStory, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)

		require.NoError(t, h.GenerateStory(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var responseBody model.Story
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseBody))
		assert.Equal(t, prompt, responseBody.Title)

		mockStoryService.AssertExpectations(t)
	})

	t.Run("faile: should return 429 Too Many Requests if limit exceeded", func(t *testing.T) {
		prompt := "A story that should fail"
		requestBody := fmt.Sprintf(`{"prompt": "%s"}`, prompt)

		mockStoryService.On("GenerateStory", testUserID, prompt).Return(nil, service.ErrGenerationLimitExceeded).Once()

		req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)

		require.NoError(t, h.GenerateStory(c))

		assert.Equal(t, http.StatusTooManyRequests, rec.Code)

		var responseBody map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseBody))
		assert.Contains(t, responseBody["error"], "daily story generation limit")

		mockStoryService.AssertExpectations(t)
	})
}

func TestStoryHandler_DeleteStory(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should delete a story", func(t *testing.T) {
		mockStoryService.On("DeleteStory", testStoryID, testUserID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		require.NoError(t, h.DeleteStory(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockStoryService.AssertExpectations(t)
	})
}

func TestStoryHandler_UpdateStory(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should update a story title", func(t *testing.T) {
		newTitle := "Updated Story Title"
		requestBody := fmt.Sprintf(`{"title": "%s"}`, newTitle)

		updatedStory := *testStory
		updatedStory.Title = newTitle

		mockStoryService.On("UpdateStoryTitle", testStoryID, testUserID, newTitle).Return(&updatedStory, nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		err := h.UpdateStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedStory model.Story
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &receivedStory))
		assert.Equal(t, newTitle, receivedStory.Title)

		mockStoryService.AssertExpectations(t)
	})
}

func TestStoryHandler_MarkStoryAsRead(t *testing.T) {
	mockStoryService, e, token := setupTestHandler(t)
	h := NewStoryHandler(mockStoryService)

	t.Run("success: should mark a story as read", func(t *testing.T) {
		mockStoryService.On("MarkStoryAsRead", testStoryID, testUserID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id/read")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		err := h.MarkStoryAsRead(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		mockStoryService.AssertExpectations(t)
	})

	t.Run("fail: should return 404 if story not found", func(t *testing.T) {
		nonExistingStoryID := 99
		mockStoryService.On("MarkStoryAsRead", nonExistingStoryID, testUserID).Return(service.ErrStoryNotFound).Once()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id/read")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(nonExistingStoryID))

		err := h.MarkStoryAsRead(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockStoryService.AssertExpectations(t)
	})
}
