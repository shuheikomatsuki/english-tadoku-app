package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- 共通セットアップ ---
func setupStoryServiceTest(t *testing.T) (*MockStoryRepository, *MockReadingRecordRepository, *MockUserRepository, *MockLLMService, IStoryService) {
	mockStoryRepo := new(MockStoryRepository)
	mockReadingRepo := new(MockReadingRecordRepository)
	mockUserRepo := new(MockUserRepository)
	mockLLM := new(MockLLMService)

	storyService := NewStoryService(mockStoryRepo, mockReadingRepo, mockUserRepo, mockLLM, testDailyLimit)

	return mockStoryRepo, mockReadingRepo, mockUserRepo, mockLLM, storyService
}

func TestStoryService_GenerateStory(t *testing.T) {
	// セットアップヘルパーを使用
	mockStoryRepo, _, mockUserRepo, mockLLM, storyService := setupStoryServiceTest(t)

	// testUser (service_test.go で定義) をコピー
	baseUser := *testUser

	t.Run("success: should generate, update count, and create story if limit not reached", func(t *testing.T) {
		prompt := "A story about mocks"
		generatedContent := "Mocks are useful for testing."
		expectedWordCount := 5

		// ユーザーの現在の状態 (0回)
		userState := baseUser
		userState.GenerationCount = 0
		userState.LastGenerationAt = nil

		// GetUserByID が呼ばれる
		mockUserRepo.On("GetUserByID", testUser.ID).Return(&userState, nil).Once()

		// LLM サービスが呼ばれる
		mockLLM.On("GenerateStory", prompt).Return(generatedContent, nil).Once()

		// StoryRepo が呼ばれる (内容は変更なし)
		mockStoryRepo.On("CreateStory", mock.AnythingOfType("*model.Story")).Return(nil).Once()

		// UpdateGenerationStatus が呼ばれる (1回に更新)
		mockUserRepo.On("UpdateGenerationStatus", testUser.ID, 1, mock.AnythingOfType("time.Time")).Return(nil).Once()

		story, err := storyService.GenerateStory(testUser.ID, prompt)

		require.NoError(t, err)
		assert.Equal(t, expectedWordCount, story.WordCount)

		mockUserRepo.AssertExpectations(t)
		mockLLM.AssertExpectations(t)
		mockStoryRepo.AssertExpectations(t)
	})

	t.Run("success: count should reset if last generation was yesterday", func(t *testing.T) {
		prompt := "A story about resetting"

		// 1. ユーザーの状態 (昨日5回生成済み)
		userState := baseUser
		userState.GenerationCount = 5
		yesterday := time.Now().AddDate(0, 0, -1) // 昨日の日付
		userState.LastGenerationAt = &yesterday

		// GetUserByID が呼ばれる
		mockUserRepo.On("GetUserByID", testUser.ID).Return(&userState, nil).Once()

		// LLM 呼び出し (カウントがリセットされ、実行される)
		mockLLM.On("GenerateStory", prompt).Return("Content", nil).Once()

		// Story 作成
		mockStoryRepo.On("CreateStory", mock.AnythingOfType("*model.Story")).Return(nil).Once()

		// カウントが 1 に更新される
		mockUserRepo.On("UpdateGenerationStatus", testUser.ID, 1, mock.AnythingOfType("time.Time")).Return(nil).Once()

		_, err := storyService.GenerateStory(testUser.ID, prompt)

		require.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail: should return ErrGenerationLimitExceeded if limit reached", func(t *testing.T) {
		prompt := "A story that should fail"

		// ユーザーの状態 (今日すでに5回生成済み)
		userState := baseUser
		userState.GenerationCount = testDailyLimit // 制限(5)に達している
		today := time.Now()
		userState.LastGenerationAt = &today

		// GetUserByID が呼ばれる
		mockUserRepo.On("GetUserByID", testUser.ID).Return(&userState, nil).Once()

		story, err := storyService.GenerateStory(testUser.ID, prompt)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrGenerationLimitExceeded)
		assert.Nil(t, story)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestStoryService_GetStories(t *testing.T) {
	// セットアップヘルパーを使用
	mockStoryRepo, _, mockUserRepo, _, storyService := setupStoryServiceTest(t)
	_ = mockUserRepo // (このテストでは使わないため、エラー回避)

	t.Run("success: should calculate pagination correctly", func(t *testing.T) {
		// ... (テストロジックは変更なし) ...
		page, limit := 2, 10
		offset := (page - 1) * limit
		totalCount := 25
		totalPages := 3
		_ = totalPages // (未使用変数エラー回避)

		mockStories := []*model.Story{testStory}
		mockStoryRepo.On("CountUserStories", testUser.ID).Return(totalCount, nil).Once()
		mockStoryRepo.On("GetUserStories", testUser.ID, limit, offset).Return(mockStories, nil).Once()
		result, err := storyService.GetStories(testUser.ID, page, limit)

		require.NoError(t, err)
		assert.Equal(t, totalCount, result.TotalCount)

		mockStoryRepo.AssertExpectations(t)
	})
}

func TestStoryService_GetStory(t *testing.T) {
	// セットアップヘルパーを使用
	mockStoryRepo, mockReadingRepo, mockUserRepo, _, storyService := setupStoryServiceTest(t)
	_ = mockUserRepo // (このテストでは使わないため、エラー回避)

	t.Run("success: should return story detail with read count", func(t *testing.T) {
		expectedReadCount := 5
		mockStoryRepo.On("GetUserStory", testStory.ID, testUser.ID).Return(testStory, nil).Once()
		mockReadingRepo.On("CountReadingRecords", testUser.ID, testStory.ID).Return(expectedReadCount, nil).Once()
		detail, err := storyService.GetStory(testStory.ID, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, expectedReadCount, detail.ReadCount)

		mockStoryRepo.AssertExpectations(t)
		mockReadingRepo.AssertExpectations(t)
	})

	t.Run("fail: should return ErrStoryNotFound if story repo returns ErrNoRows", func(t *testing.T) {
		mockStoryRepo.On("GetUserStory", testStory.ID, testUser.ID).Return(nil, sql.ErrNoRows).Once()
		_, err := storyService.GetStory(testStory.ID, testUser.ID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrStoryNotFound)

		mockStoryRepo.AssertExpectations(t)
		mockReadingRepo.AssertExpectations(t)
	})
}

func TestStoryService_MarkStoryAsRead(t *testing.T) {
	// セットアップヘルパーを使用
	mockStoryRepo, mockReadingRepo, mockUserRepo, _, storyService := setupStoryServiceTest(t)
	_ = mockUserRepo // (このテストでは使わないため、Linterエラー回避)

	t.Run("success: should create reading record with correct word count", func(t *testing.T) {
		mockStoryRepo.On("GetUserStory", testStory.ID, testUser.ID).Return(testStory, nil).Once()
		mockReadingRepo.On("CreateReadingRecord", testUser.ID, testStory.ID, testStory.WordCount).Return(nil).Once()
		err := storyService.MarkStoryAsRead(testStory.ID, testUser.ID)
		require.NoError(t, err)

		mockStoryRepo.AssertExpectations(t)
		mockReadingRepo.AssertExpectations(t)
	})
}
