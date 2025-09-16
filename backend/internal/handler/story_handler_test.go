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

const (
	testUserID = 123
	testStoryID = 1
)

// テスト用のストーリーオブジェクトを生成するヘルパー関数
func newTestStory(id, userID int) *model.Story {
	return &model.Story{
		ID: id,
		UserID: userID,
		Title: fmt.Sprintf("Story %d", id),
		Content: fmt.Sprintf("This is the content for story %d", id),
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

	t.Run("success: should return a story", func(t *testing.T) {
		expectedStory := &model.Story{
			ID: 1,
			UserID: 123,
			Title: "Test Story",
			Content: "This is a test.",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetUserStory", 1, 123).Return(expectedStory, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		// TODO: JWTミドルウェア認証を実装したら、テスト用のユーザー情報を Context にセットする
		// claims := &tokenClaims{UserID: 123}
		// c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))

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

	t.Run("fail: should return 404 Not Found for non-existent story", func(t *testing.T) {
		nonExistingStoryID := 99

		mockRepo.On("GetUserStory", nonExistingStoryID, testUserID).Return(nil, http.ErrMissingFile).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(nonExistingStoryID))

		// TODO: JWTミドルウェア認証を実装したら、テスト用のユーザー情報を Context にセットする
		// claims := &tokenClaims{UserID: 123}
		// c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))

		err := h.GetStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestStoryHandler_GetStories(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

	t.Run("success: should return a list of story", func(t *testing.T) {
		expectedStories := []*model.Story{
			newTestStory(1, testUserID),
			newTestStory(2, testUserID),
		}

		limit, offset := 10, 0
		// TODO: limit, offsetをクエリパラメータから取得するようになったら、それらも引数に含める
		mockRepo.On("GetUserStories", testUserID, limit, offset).Return(expectedStories, nil).Once()

		reqURL := fmt.Sprintf("/stories?limit=%d&offset=%d", limit, offset)
		req := httptest.NewRequest(http.MethodGet, reqURL, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// TODO: JWTミドルウェアを実装したら、ユーザー情報を Context にセットする。
		// c.Set("user", ...)

		err := h.GetStories(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedStories []model.Story
		err = json.Unmarshal(rec.Body.Bytes(), &receivedStories)
		require.NoError(t, err)
		assert.Len(t, receivedStories, len(expectedStories))
		assert.Equal(t, expectedStories[0].Title, receivedStories[0].Title)

		mockRepo.AssertExpectations(t)
	})
}

func TestStoryHandler_GenerateStory(t *testing.T) {
	mockRepo := new(MockStoryRepository)
	mockLLM := new(MockLLMService)
	h := NewStoryHandler(mockRepo, mockLLM)
	e := echo.New()

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
	
		// TODO: User Context をセットする
	
		err := h.GenerateStory(c)
	
		require.NoError(t, err)
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

	t.Run("success: should delete a story", func(t *testing.T) {
		mockRepo.On("DeleteStory", testStoryID).Return(nil).Once()
		// TODO: 削除権限チェックが必要のため、先に story を取得する必要があるかも
		// mockRepo.On("GetUserStory", 1, 123).Return(&model.Story{ID: 1, UserID: 123}, nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(testStoryID))

		// TODO: JWTミドルウェアからユーザー情報を Context にセットする
		// c.Set("user", ...)

		err := h.DeleteStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("fail: should return 400 Bad Request for invalid id", func(t *testing.T) {
		// mockRepo.ON(...)を書かないのは、モックのメソッドが呼ばれないから。

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/stories/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid_id")

		err := h.DeleteStory(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}