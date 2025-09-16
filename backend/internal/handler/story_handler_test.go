package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	// "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type MockStoryRepository struct {
	mock.Mock
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

type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) GenerateStory(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

// テストケースを記述
func TestStoryHandler_GetStory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockStoryRepository)
		mockLLM := new(MockLLMService)

		expectedStory := &model.Story{
			ID: 1,
			UserID: 123,
			Title: "Test Story",
			Content: "This is a test.",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetUserStory", 1, 123).Return(expectedStory, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// TODO: JWTミドルウェア認証を実装したら、テスト用のユーザー情報を Context にセットする
		// claims := &tokenClaims{UserID: 123}
		// c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))

		h := NewStoryHandler(mockRepo, mockLLM)

		err := h.GetStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedStory model.Story
		err = json.Unmarshal(rec.Body.Bytes(), &receivedStory)
		require.NoError(t, err)
		assert.Equal(t, expectedStory.Title, receivedStory.Title)
		assert.Equal(t, expectedStory.Content, receivedStory.Content)

		mockRepo.AssertExpectations(t)
	})

	t.Run("404 Not Found", func(t *testing.T) {
		mockRepo := new(MockStoryRepository)
		mockLLM := new(MockLLMService)

		mockRepo.On("GetUserStory", 99, 123).Return(nil, http.ErrMissingFile)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues("99")

		// TODO: JWTミドルウェア認証を実装したら、テスト用のユーザー情報を Context にセットする
		// claims := &tokenClaims{UserID: 123}
		// c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))

		h := NewStoryHandler(mockRepo, mockLLM)

		err := h.GetStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestStoryHandler_GenerateStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)

	requestBody := `{"prompt": "What is Golang?"}`

	mockLLM.On("GenerateStory", "What is Golang?").Return("Golang is an open-source programming language developed by Google.", nil)
	mockRepo.On("CreateStory", mock.AnythingOfType("*model.Story")).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/stories", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// TODO: User Context をセットする

	h := NewStoryHandler(mockRepo, mockLLM)

	err := h.GenerateStory(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	assert.Contains(t, responseBody["content"], "Golang is")

	mockRepo.AssertExpectations(t)
	mockLLM.AssertExpectations(t)
}