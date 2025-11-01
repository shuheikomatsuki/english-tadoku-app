package handler

import (
	"database/sql"
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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type MockStoryRepository struct {
	mock.Mock
}

func (m *MockStoryRepository) CountUserStories(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockStoryRepository) CreateStory(story *model.Story) error {
	args := m.Called(story)
	if args.Error(0) == nil {
		story.ID = 1 // 仮のIDを設定
	}
	return args.Error(0)
}

func (m *MockStoryRepository) GetUserStories(userID, limit, offset int) ([]*model.Story, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Story), args.Error(1)
}

func (m *MockStoryRepository) GetUserStory(id, userID int) (*model.Story, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

func (m *MockStoryRepository) DeleteStory(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStoryRepository) UpdateStoryTitle(storyID, userID int, newTitle string) (*model.Story, error) {
	args := m.Called(storyID, userID, newTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

func (m *MockStoryRepository) CreateReadingRecord(userID, storyID, wordCount int) error {
	args := m.Called(userID, storyID, wordCount)
	return args.Error(0)
}

func (m *MockStoryRepository) CountReadingRecords(userID, storyID int) (int, error)  {
	args := m.Called(userID, storyID)
	return args.Int(0), args.Error(1)
}

func (m *MockStoryRepository) GetLatestReadingRecord(userID, storyID int) (*model.ReadingRecord, error) {
	args := m.Called(userID, storyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReadingRecord), args.Error(1)
}

func (m *MockStoryRepository) DeleteReadingRecord(recordID int, userID int) error {
	args := m.Called(recordID, userID)
	return args.Error(0)
}

type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) GenerateStory(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

const (
	testUserID  = 123
	testStoryID = 1
)

// テスト用のストーリーオブジェクトを生成するヘルパー関数
func newTestStory(id, userID int) *model.Story {
	return &model.Story{
		ID:        id,
		UserID:    userID,
		Title:     fmt.Sprintf("Story %d", id),
		Content:   fmt.Sprintf("This is the content for story %d", id),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// --- テストケース ---
func TestStoryHandler_GetStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should return a story", func(t *testing.T) {
		expectedStory := newTestStory(testStoryID, testUserID)
		mockRepo.On("GetUserStory", testStoryID, testUserID).Return(expectedStory, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		require.NoError(t, h.GetStory(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedStory model.Story
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &receivedStory))
		assert.Equal(t, expectedStory.Title, receivedStory.Title)
		assert.Equal(t, expectedStory.Content, receivedStory.Content)

		mockRepo.AssertExpectations(t)
	})

	t.Run("fail: should return 404 Not Found for non-existent story", func(t *testing.T) {
		nonExistingStoryID := 99
		mockRepo.On("GetUserStory", nonExistingStoryID, testUserID).Return(nil, http.ErrMissingFile).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(nonExistingStoryID))

		require.NoError(t, h.GetStory(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestStoryHandler_GetStories(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should return a list of story", func(t *testing.T) {
		expectedStories := []*model.Story{
			newTestStory(1, testUserID),
			newTestStory(2, testUserID),
		}

		// limit, offset := 10, 0
		page, limit := 1, 10
		offset := (page - 1) * limit

		mockRepo.On("CountUserStories", testUserID).Return(len(expectedStories), nil).Once()

		mockRepo.On("GetUserStories", testUserID, limit, offset).Return(expectedStories, nil).Once()

		reqURL := fmt.Sprintf("/stories?limit=%d&offset=%d", limit, offset)
		req := httptest.NewRequest(http.MethodGet, reqURL, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)

		require.NoError(t, h.GetStories(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		// var receivedStories []model.Story
		// require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &receivedStories))
		// assert.Len(t, receivedStories, len(expectedStories))
		// assert.Equal(t, expectedStories[0].Title, receivedStories[0].Title)

		var response GetStoriesResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Len(t, response.Stories, len(expectedStories))
		assert.Equal(t, expectedStories[0].Title, response.Stories[0].Title)
		assert.Equal(t, len(expectedStories), response.TotalCount)
		assert.Equal(t, 1, response.TotalPages)
		assert.Equal(t, page, response.CurrentPage)

		mockRepo.AssertExpectations(t)
	})
}

func TestStoryHandler_GenerateStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should generate and save a story", func(t *testing.T) {
		prompt := "A story abount Go"
		generatedContent := "Go is a statically typed, compiled programming language..."
		requestBody := fmt.Sprintf(`{"prompt": "%s"}`, prompt)

		mockLLM.On("GenerateStory", prompt).Return(generatedContent, nil).Once()
		mockRepo.On("CreateStory", mock.AnythingOfType("*model.Story")).Return(nil).Once()

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
		assert.Equal(t, generatedContent, responseBody.Content)
		mockRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
	})
}

func TestStoryHandler_DeleteStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should delete a story", func(t *testing.T) {
		mockRepo.On("GetUserStory", testStoryID, testUserID).Return(&model.Story{ID: testStoryID, UserID: testUserID}, nil).Once()
		mockRepo.On("DeleteStory", testStoryID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		require.NoError(t, h.DeleteStory(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail: should return 400 Bad Request for invalid id", func(t *testing.T) {
		// mockRepo.ON(...)を書かないのは、モックのメソッドが呼ばれないから。

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid_id")

		err := h.DeleteStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStoryHandler_UpdateStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()
	e.Validator = NewValidator()

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should update a story title", func(t *testing.T) {
		newTitle := "Updated Story Title"
		requestBody := fmt.Sprintf(`{"title": "%s"}`, newTitle)

		updatedStory := newTestStory(testStoryID, testUserID)
		updatedStory.Title = newTitle

		mockRepo.On("UpdateStoryTitle", testStoryID, testUserID, newTitle).Return(updatedStory, nil).Once()

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
		assert.Equal(t, testStoryID, receivedStory.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("fail: should return 400 for invalid request body", func(t *testing.T) {
		requestBody := `{"title": ""}`

		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		err := h.UpdateStory(c)

		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})
}

func TestStoryHandler_MarkStoryAsRead(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	claims := &JwtCustomClaims{
		testUserID, 
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should mark a story as read", func(t *testing.T) {
		storyToRead := newTestStory(testStoryID, testUserID)
		storyToRead.WordCount = 50

		mockRepo.On("GetUserStory", testStoryID, testUserID).Return(storyToRead, nil).Once()
		mockRepo.On("CreateReadingRecord", testUserID, testStoryID, storyToRead.WordCount).Return(nil).Once()

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
		
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail: should return 404 if story not found", func(t *testing.T) {
		nonExistingStoryID := 99
		mockRepo.On("GetUserStory", nonExistingStoryID, testUserID).Return(nil, sql.ErrNoRows).Once()

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

		mockRepo.AssertExpectations(t)
	})
}