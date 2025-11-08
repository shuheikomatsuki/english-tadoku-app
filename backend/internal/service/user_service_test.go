package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// テスト用の制限値を定義
const testDailyLimit = 5

func TestUserService_GetUserStats(t *testing.T) {
	mockReadingRepo := new(MockReadingRecordRepository)
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockReadingRepo, mockUserRepo, testDailyLimit)

	t.Run("success: should calculate all stats correctly", func(t *testing.T) {
		mockReadingRepo.On("GetWordCountInDateRange", testUser.ID, mock.Anything, mock.Anything).
			Return(100, nil). // 1回目 (Today)
			Once()
		mockReadingRepo.On("GetWordCountInDateRange", testUser.ID, mock.Anything, mock.Anything).
			Return(700, nil). // 2回目 (Weekly)
			Once()
		mockReadingRepo.On("GetWordCountInDateRange", testUser.ID, mock.Anything, mock.Anything).
			Return(3000, nil). // 3回目 (Monthly)
			Once()
		mockReadingRepo.On("GetWordCountInDateRange", testUser.ID, mock.Anything, mock.Anything).
			Return(10000, nil). // 4回目 (Yearly)
			Once()

		mockReadingRepo.On("GetUserTotalWordCount", testUser.ID).Return(50000, nil).Once()
		now := time.Now()
		todayKey := now.Format("2006-01-02")
		mockLast7Days := map[string]int{todayKey: 100}
		mockReadingRepo.On("GetDailyWordCountLastNDays", testUser.ID, 7, mock.AnythingOfType("time.Time")).Return(mockLast7Days, nil).Once()

		stats, err := userService.GetUserStats(testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, 100, stats.TodayWordCount)
		assert.Equal(t, 700, stats.WeeklyWordCount)
		assert.Equal(t, 3000, stats.MonthlyWordCount)
		assert.Equal(t, 10000, stats.YearlyWordCount)
		assert.Equal(t, 50000, stats.TotalWordCount)
		assert.Equal(t, 100, stats.Last7DaysWordCount[todayKey])
		
		mockReadingRepo.AssertExpectations(t)
	})
}

func TestUserService_GetGenerationStatus(t *testing.T) {
	mockReadingRepo := new(MockReadingRecordRepository)
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockReadingRepo, mockUserRepo, testDailyLimit)

	baseUser := *testUser

	t.Run("success: should return current count if generated today", func(t *testing.T) {
		userState := baseUser
		userState.GenerationCount = 3
		today := time.Now()
		userState.LastGenerationAt = &today

		mockUserRepo.On("GetUserByID", testUser.ID).Return(&userState, nil).Once()

		status, err := userService.GetGenerationStatus(testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, 3, status.CurrentCount)
		assert.Equal(t, testDailyLimit, status.Limit)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success: should return 0 if last generation was yesterday", func(t *testing.T) {
		userState := baseUser
		userState.GenerationCount = 5
		yesterday := time.Now().AddDate(0, 0, -1)
		userState.LastGenerationAt = &yesterday

		mockUserRepo.On("GetUserByID", testUser.ID).Return(&userState, nil).Once()

		status, err := userService.GetGenerationStatus(testUser.ID)

		require.NoError(t, err)
		assert.Equal(t, 0, status.CurrentCount)
		assert.Equal(t, testDailyLimit, status.Limit)

		mockUserRepo.AssertExpectations(t)
	})
}