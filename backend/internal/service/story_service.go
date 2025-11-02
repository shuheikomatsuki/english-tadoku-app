package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
)

// const (
// 	dailyGenerationLimit = 1
// )

// PaginatedStories はサービス層が返すページネーション結果のモデル
type PaginatedStories struct {
	Stories     []*model.Story
	TotalCount  int
	TotalPages  int
	CurrentPage int
}

// StoryDetail はサービス層が返すストーリー詳細のモデル
type StoryDetail struct {
	model.Story
	ReadCount int
}

// サービス層で扱うためのドメインエラーを定義
var (
	ErrStoryNotFound           = errors.New("story not found")
	ErrForbidden               = errors.New("forbidden")
	ErrNoReadingRecord         = errors.New("no reading record found")
	ErrGenerationLimitExceeded = errors.New("generation limit exceeded")
)

type IStoryService interface {
	GenerateStory(userID int, prompt string) (*model.Story, error)
	GetStories(userID int, page, limit int) (*PaginatedStories, error)
	GetStory(storyID, userID int) (*StoryDetail, error)
	DeleteStory(storyID, userID int) error
	UpdateStoryTitle(storyID, userID int, newTitle string) (*model.Story, error)
	MarkStoryAsRead(storyID, userID int) error
	UndoLastRead(storyID, userID int) error
}

type StoryService struct {
	StoryRepo         repository.IStoryRepository
	ReadingRecordRepo repository.IReadingRecordRepository
	UserRepo          repository.IUserRepository
	LLMService        ILLMService // llm_service.go に依存
}

func NewStoryService(storyRepo repository.IStoryRepository, readingRecordRepo repository.IReadingRecordRepository, userRepo repository.IUserRepository, llmService ILLMService) IStoryService {
	return &StoryService{
		StoryRepo:         storyRepo,
		ReadingRecordRepo: readingRecordRepo,
		UserRepo:          userRepo,
		LLMService:        llmService,
	}
}

func (s *StoryService) GenerateStory(userID int, prompt string) (*model.Story, error) {

	// ユーザーの生成制限を確認
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for validation: %w", err)
	}

	now := time.Now()
	todayStart := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)

	currentCount := user.GenerationCount
	lastGen := user.LastGenerationAt

	if lastGen != nil && lastGen.UTC().Before(todayStart) {
		currentCount = 0
	}

	if currentCount >= dailyGenerationLimit {
		return nil, ErrGenerationLimitExceeded
	}

	// LLMサービス呼び出し
	content, err := s.LLMService.GenerateStory(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate story: %w", err)
	}

	wordCount := len(strings.Fields(content))

	story := &model.Story{
		UserID:    userID,
		Title:     prompt,
		Content:   content,
		WordCount: wordCount,
	}

	// DB保存
	if err := s.StoryRepo.CreateStory(story); err != nil {
		return nil, fmt.Errorf("failed to save story: %w", err)
	}

	newCount := currentCount + 1
	if err := s.UserRepo.UpdateGenerationStatus(userID, newCount, now); err != nil {
		log.Printf("WARNING: failed to update generation status for user %d: %v", userID, err)
	}

	return story, nil
}

func (s *StoryService) GetStories(userID int, page, limit int) (*PaginatedStories, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	totalCount, err := s.StoryRepo.CountUserStories(userID)
	if err != nil {
		return nil, fmt.Errorf("database error (count): %w", err)
	}

	stories, err := s.StoryRepo.GetUserStories(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error (get): %w", err)
	}

	if stories == nil {
		stories = []*model.Story{}
	}

	// ページネーション計算
	totalPages := 0
	if totalCount > 0 {
		totalPages = int(math.Ceil(float64(totalCount) / float64(limit)))
	}

	res := &PaginatedStories{
		Stories:     stories,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	return res, nil
}

func (s *StoryService) GetStory(storyID, userID int) (*StoryDetail, error) {
	story, err := s.StoryRepo.GetUserStory(storyID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStoryNotFound
		}
		return nil, fmt.Errorf("database error (get story): %w", err)
	}

	readCount, err := s.ReadingRecordRepo.CountReadingRecords(userID, storyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reading count: %w", err)
	}

	res := &StoryDetail{
		Story:     *story,
		ReadCount: readCount,
	}

	return res, nil
}

func (s *StoryService) checkStoryOwnership(storyID, userID int) (*model.Story, error) {
	story, err := s.StoryRepo.GetUserStory(storyID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStoryNotFound
		}
		return nil, fmt.Errorf("database error (get story): %w", err)
	}

	if story.UserID != userID {
		return nil, ErrForbidden
	}
	return story, nil
}

func (s *StoryService) DeleteStory(storyID, userID int) error {
	_, err := s.checkStoryOwnership(storyID, userID)
	if err != nil {
		return err
	}

	if err := s.StoryRepo.DeleteStory(storyID); err != nil {
		return fmt.Errorf("failed to delete story: %w", err)
	}

	return nil
}

func (s *StoryService) UpdateStoryTitle(storyID, userID int, newTitle string) (*model.Story, error) {
	updatedStory, err := s.StoryRepo.UpdateStoryTitle(storyID, userID, newTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStoryNotFound
		}
		return nil, fmt.Errorf("failed to update story: %w", err)
	}

	return updatedStory, nil
}

func (s *StoryService) MarkStoryAsRead(storyID, userID int) error {
	// ストーリーの存在と所有権を確認
	story, err := s.checkStoryOwnership(storyID, userID)
	if err != nil {
		return err
	}
	// 読書記録を作成
	err = s.ReadingRecordRepo.CreateReadingRecord(userID, story.ID, story.WordCount)
	if err != nil {
		return fmt.Errorf("failed to create reading record: %w", err)
	}
	return nil
}

func (s *StoryService) UndoLastRead(storyID, userID int) error {
	// 最新の読書記録を取得
	latestRecord, err := s.ReadingRecordRepo.GetLatestReadingRecord(userID, storyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoReadingRecord
		}
		return fmt.Errorf("database error (get latest record): %w", err)
	}

	// 読書記録を削除
	err = s.ReadingRecordRepo.DeleteReadingRecord(latestRecord.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reading record: %w", err)
	}

	return nil
}
