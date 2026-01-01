package service

import (
	"fmt"
	"time"

	"github.com/shuheikomatsuki/readoku/backend/internal/repository"
)

// const (
// 	dailyGenerationLimit = 5
// )

type UserStats struct {
	TotalWordCount     int
	TodayWordCount     int
	WeeklyWordCount    int
	MonthlyWordCount   int
	YearlyWordCount    int
	Last7DaysWordCount map[string]int
}

type GenerationStatus struct {
	CurrentCount int `json:"current_count"`
	Limit        int `json:"limit"`
}

type IUserService interface {
	GetUserStats(userID int) (*UserStats, error)
	GetGenerationStatus(userID int) (*GenerationStatus, error)
}

type UserService struct {
	ReadingRecordRepo repository.IReadingRecordRepository
	UserRepo          repository.IUserRepository
	DailyLimit        int
}

func NewUserService(readingRecordRepo repository.IReadingRecordRepository, userRepo repository.IUserRepository, dailyLimit int) IUserService {
	return &UserService{
		ReadingRecordRepo: readingRecordRepo,
		UserRepo:          userRepo,
		DailyLimit:        dailyLimit,
	}
}

func (s *UserService) GetUserStats(userID int) (*UserStats, error) {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayCount, err := s.ReadingRecordRepo.GetWordCountInDateRange(userID, startOfDay, now.Add(24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (today): %w", err)
	}

	// 週の始まりを月曜日とする
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜日の場合は7に調整
	}
	startOfWeek := startOfDay.AddDate(0, 0, -weekday+1)
	weeklyCount, err := s.ReadingRecordRepo.GetWordCountInDateRange(userID, startOfWeek, startOfWeek.AddDate(0, 0, 7))
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (weekly): %w", err)
	}

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthlyCount, err := s.ReadingRecordRepo.GetWordCountInDateRange(userID, startOfMonth, startOfMonth.AddDate(0, 1, 0))
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (monthly): %w", err)
	}

	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	yearlyCount, err := s.ReadingRecordRepo.GetWordCountInDateRange(userID, startOfYear, startOfYear.AddDate(1, 0, 0))
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (yearly): %w", err)
	}

	totalWordCount, err := s.ReadingRecordRepo.GetUserTotalWordCount(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (total): %w", err)
	}

	last7DaysCount, err := s.ReadingRecordRepo.GetDailyWordCountLastNDays(userID, 7, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats (last7days): %w", err)
	}

	if last7DaysCount == nil {
		last7DaysCount = make(map[string]int)
	}

	stats := &UserStats{
		TotalWordCount:     totalWordCount,
		TodayWordCount:     todayCount,
		WeeklyWordCount:    weeklyCount,
		MonthlyWordCount:   monthlyCount,
		YearlyWordCount:    yearlyCount,
		Last7DaysWordCount: last7DaysCount,
	}

	return stats, nil
}

func (s *UserService) GetGenerationStatus(userID int) (*GenerationStatus, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()
	todayStart := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)

	currentCount := user.GenerationCount
	lastGen := user.LastGenerationAt

	if lastGen != nil && lastGen.UTC().Before(todayStart) {
		currentCount = 0
	}

	// const dailyLimit = 5

	status := &GenerationStatus{
		CurrentCount: currentCount,
		Limit:        s.DailyLimit,
	}

	return status, nil
}
