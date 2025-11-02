package handler

import (
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}

func (m *MockAuthService) ValidateUser(email, password string) (*model.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) GenerateToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserStats(userID int) (*service.UserStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserStats), args.Error(1)
}

type MockStoryService struct {
	mock.Mock
}

func (m *MockStoryService) GenerateStory(userID int, prompt string) (*model.Story, error) {
	args := m.Called(userID, prompt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

func (m *MockStoryService) GetStories(userID int, page, limit int) (*service.PaginatedStories, error) {
	args := m.Called(userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.PaginatedStories), args.Error(1)
}

func (m *MockStoryService) GetStory(storyID, userID int) (*service.StoryDetail, error) {
	args := m.Called(storyID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.StoryDetail), args.Error(1)
}

func (m *MockStoryService) DeleteStory(storyID, userID int) error {
	args := m.Called(storyID, userID)
	return args.Error(0)
}

func (m *MockStoryService) UpdateStoryTitle(storyID, userID int, newTitle string) (*model.Story, error) {
	args := m.Called(storyID, userID, newTitle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Story), args.Error(1)
}

func (m *MockStoryService) MarkStoryAsRead(storyID, userID int) error {
	args := m.Called(storyID, userID)
	return args.Error(0)
}

func (m *MockStoryService) UndoLastRead(storyID, userID int) error {
	args := m.Called(storyID, userID)
	return args.Error(0)
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