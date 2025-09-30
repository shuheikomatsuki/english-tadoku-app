package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type IStoryRepository interface {
	CreateStory(story *model.Story) error
	GetUserStories(userID, limit, offset int) ([]*model.Story, error)
	CountUserStories(userID int) (int, error)
	GetUserStory(storyID int, userID int) (*model.Story, error)
	DeleteStory(storyID int) error
	UpdateStoryTitle(storyID int, userID int, newTitle string) (*model.Story, error)
}

type sqlxStoryRepository struct {
	DB *sqlx.DB
}

func NewStoryRepository(db *sqlx.DB) IStoryRepository {
	return &sqlxStoryRepository{DB: db}
}

func (r *sqlxStoryRepository) CreateStory(story *model.Story) error {
	query := `
		INSERT INTO stories(user_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRowx(query, story.UserID, story.Title, story.Content).Scan(&story.ID, &story.CreatedAt, &story.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create story: %w", err)
	}
	return nil
}

func (r *sqlxStoryRepository) GetUserStories(userID, limit, offset int) ([]*model.Story, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM stories
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var stories []*model.Story
	err := r.DB.Select(&stories, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stories: %w", err)
	}
	return stories, nil
}

func (r *sqlxStoryRepository) CountUserStories(userID int) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM stories WHERE user_id = $1`
	err := r.DB.Get(&total, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count user stories: %w", err)
	}
	return total, nil
}

func (r *sqlxStoryRepository) GetUserStory(storyID int, userID int) (*model.Story, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM stories
		WHERE id = $1 AND user_id = $2
	`
	var story model.Story
	err := r.DB.Get(&story, query, storyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user story: %w", err)
	}
	return &story, nil
}

func (r *sqlxStoryRepository) DeleteStory(storyID int) error {
	query := `
		DELETE FROM stories
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, storyID)
	if err != nil {
		return fmt.Errorf("failed to delete story: %w", err)
	}
	return nil
}

func (r *sqlxStoryRepository) UpdateStoryTitle(storyID int, userID int, newTitle string) (*model.Story, error) {
	var updatedStory model.Story
	query := `
		UPDATE stories
		SET title = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, title, content, created_at, updated_at
	`
	err := r.DB.Get(&updatedStory, query, newTitle, storyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update story: %w", err)
	}
	return &updatedStory, nil
}