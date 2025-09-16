package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
}
type IStoryRepository interface {
	CreateStory(story *model.Story) error
	GetUserStories(userID, limit, offset int) ([]*model.Story, error) // ユーザーのストーリー一覧を取得（タイトルのみ）
	GetUserStory(storyID int, userID int) (*model.Story, error)       // ユーザーの特定のストーリーを取得（タイトルと内容）
	DeleteStory(storyID int) error
}

type UserRepository struct {
	DB *sqlx.DB
}

type StoryRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRowx(query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.DB.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func NewStoryRepository(db *sqlx.DB) IStoryRepository {
	return &StoryRepository{DB: db}
}

func (r *StoryRepository) CreateStory(story *model.Story) error {
	// 作成した story をDBに保存する処理
	query := `
		INSERT INTO stories(user_id, title, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.DB.QueryRowx(query, story.UserID, story.Title, story.Content).Scan(&story.ID)
	if err != nil {
		return fmt.Errorf("failed to create story: %w", err)
	}

	return nil
}

func (r *StoryRepository) GetUserStories(userID, limit, offset int) ([]*model.Story, error) {
	// 指定された userID の story 一覧をDBから取得する処理（タイトルのみ）
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

func (r *StoryRepository) GetUserStory(storyID int, userID int) (*model.Story, error) {
	// 指定された storyID と userID の story をDBから取得する処理（タイトルと内容）
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

func (r *StoryRepository) DeleteStory(storyID int) error {
	// 指定された storyID の story をDBから削除する処理
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
