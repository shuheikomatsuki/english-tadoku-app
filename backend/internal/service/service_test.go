package service

import (
	"time"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.User), args.Error(1)
}

type MockStoryRepository struct {
	mock.Mock
}

func (m *MockStoryRepository) CreateStory(story *model.Story) error {
	args := m.Called(story)
	if args.Error(0) == nil {
		story.ID = 1 // モックでIDを設定
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

func (m *MockStoryRepository) CountUserStories(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockStoryRepository) GetUserStory(storyID, userID int) (*model.Story, error) {
	args := m.Called(storyID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

func (m *MockStoryRepository) DeleteStory(storyID int) error {
	args := m.Called(storyID)
	return args.Error(0)
}

func (m *MockStoryRepository) UpdateStoryTitle(storyID, userID int, newTitle string) (*model.Story, error) {
	args := m.Called(storyID, userID, newTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

type MockReadingRecordRepository struct {
	mock.Mock
}

func (m *MockReadingRecordRepository) CreateReadingRecord(userID, storyID, wordCount int) error {
	args := m.Called(userID, storyID, wordCount)
	return args.Error(0)
}

func (m *MockReadingRecordRepository) CountReadingRecords(userID, storyID int) (int, error) {
	args := m.Called(userID, storyID)
	return args.Int(0), args.Error(1)
}

func (m *MockReadingRecordRepository) GetLatestReadingRecord(userID, storyID int) (*model.ReadingRecord, error) {
	args := m.Called(userID, storyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReadingRecord), args.Error(1)
}

func (m *MockReadingRecordRepository) DeleteReadingRecord(recordID int, userID int) error {
	args := m.Called(recordID, userID)
	return args.Error(0)
}

func (m *MockReadingRecordRepository) GetUserTotalWordCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockReadingRecordRepository) GetWordCountInDateRange(userID int, start, end time.Time) (int, error) {
	// time.Time は mock.Anything を使うことが多いらしい
	args := m.Called(userID, start, end)
	// mock.Anything を使う場合
	// args := m.Called(userID, mock.Anything, mock.Anything)
	return args.Int(0), args.Error(1)
}

func (m *MockReadingRecordRepository) GetDailyWordCountLastNDays(userID, days int, anchorTime time.Time) (map[string]int, error) {
	args := m.Called(userID, days, anchorTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) GenerateStory(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

var testUser = &model.User{
	ID:           1,
	Email:        "test@example.com",
	PasswordHash: "$2a$10$Q4A86sZk6FTXTZTonDVMm.npTH5yp2e8/vmvk2EWKLOGxmaPF127a",
}

var testStory = &model.Story{
	ID:        10,
	UserID:    testUser.ID,
	Title:     "Test Story",
	Content:   "This is a test story content.",
	WordCount: 6,
}
