package service

import (
	"database/sql"
	"testing"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStoryService_GenerateStory(t *testing.T) {
	mockStoryRepo := new(MockStoryRepository)
	mockReadingRepo := new(MockReadingRecordRepository)
	mockLLM := new(MockLLMService)
	storyService := NewStoryService(mockStoryRepo, mockReadingRepo, mockLLM)

	t.Run("success: should generate content, count words, and create story", func(t *testing.T) {
		prompt := "A story about mocks"
		generatedContent := "Mocks are useful for testing."
		expectedWordCount := 5

		mockLLM.On("GenerateStory", prompt).Return(generatedContent, nil).Once()

		mockStoryRepo.On("CreateStory", mock.AnythingOfType("*model.Story")).Run(func(args mock.Arguments) {
			storyArg := args.Get(0).(*model.Story)
			assert.Equal(t, expectedWordCount, storyArg.WordCount)
			assert.Equal(t, prompt, storyArg.Title)
			assert.Equal(t, generatedContent, storyArg.Content)
			assert.Equal(t, testUser.ID, storyArg.UserID)
		}).Return(nil).Once()

		story, err := storyService.GenerateStory(testUser.ID, prompt)

		require.NoError(t, err)
		assert.Equal(t, expectedWordCount, story.WordCount)

		mockLLM.AssertExpectations(t)
		mockStoryRepo.AssertExpectations(t)
	})
}

func TestStoryService_GetStories(t *testing.T) {
	mockStoryRepo := new(MockStoryRepository)
	mockReadingRepo := new(MockReadingRecordRepository)
	mockLLM := new(MockLLMService)
	storyService := NewStoryService(mockStoryRepo, mockReadingRepo, mockLLM)

	t.Run("success: should calculate pagination correctly", func(t *testing.T) {
		page, limit := 2, 10
		offset := (page - 1) * limit // 10
		totalCount := 25
		totalPages := 3

		mockStories := []*model.Story{testStory}

		mockStoryRepo.On("CountUserStories", testUser.ID).Return(totalCount, nil).Once()
		mockStoryRepo.On("GetUserStories", testUser.ID, limit, offset).Return(mockStories, nil).Once()

		result, err := storyService.GetStories(testUser.ID, page, limit)

		require.NoError(t, err)
		assert.Equal(t, totalCount, result.TotalCount)
		assert.Equal(t, totalPages, result.TotalPages)
		assert.Equal(t, page, result.CurrentPage)
		assert.Equal(t, testStory.Title, result.Stories[0].Title)

		mockStoryRepo.AssertExpectations(t)
	})
}

func TestStoryService_GetStory(t *testing.T) {
	mockStoryRepo := new(MockStoryRepository)
	mockReadingRepo := new(MockReadingRecordRepository)
	mockLLM := new(MockLLMService)
	storyService := NewStoryService(mockStoryRepo, mockReadingRepo, mockLLM)

	t.Run("success: should return story detail with read count", func(t *testing.T) {
		expectedReadCount := 5

		mockStoryRepo.On("GetUserStory", testStory.ID, testUser.ID).Return(testStory, nil).Once()
		mockReadingRepo.On("CountReadingRecords", testUser.ID, testStory.ID).Return(expectedReadCount, nil).Once()

		detail, err := storyService.GetStory(testStory.ID, testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, testStory.Title, detail.Story.Title)
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
	mockStoryRepo := new(MockStoryRepository)
	mockReadingRepo := new(MockReadingRecordRepository)
	mockLLM := new(MockLLMService)
	storyService := NewStoryService(mockStoryRepo, mockReadingRepo, mockLLM)

	t.Run("success: should create reading record with correct word count", func(t *testing.T) {
		mockStoryRepo.On("GetUserStory", testStory.ID, testUser.ID).Return(testStory, nil).Once()

		mockReadingRepo.On("CreateReadingRecord", testUser.ID, testStory.ID, testStory.WordCount).Return(nil).Once()

		err := storyService.MarkStoryAsRead(testStory.ID, testUser.ID)

		require.NoError(t, err)

		mockStoryRepo.AssertExpectations(t)
		mockReadingRepo.AssertExpectations(t)
	})
}