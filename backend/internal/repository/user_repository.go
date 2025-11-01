package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type IUserRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	GetUserTotalWordCount(userID int) (int, error)
	GetWordCountInDateRange(userID int, start, end time.Time) (int, error)
	GetDailyWordCountLastNDays(userID, days int) (map[string]int, error)
}

type sqlxUserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &sqlxUserRepository{DB: db}
}

func (r *sqlxUserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRowx(query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *sqlxUserRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.DB.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *sqlxUserRepository) GetUserTotalWordCount(userID int) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(word_count), 0) FROM reading_records WHERE user_id = $1`
	
	err := r.DB.Get(&total, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user total word count: %w", err)
	}

	return total, nil
}

func (r *sqlxUserRepository) GetWordCountInDateRange(userID int, start, end time.Time) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(word_count), 0)
		FROM reading_records
		WHERE user_id = $1 AND read_at >= $2 AND read_at < $3
	`
	err := r.DB.Get(&total, query, userID, start, end)
	if err != nil {
		return 0, fmt.Errorf("failed to get word count in date range: %w", err)
	}
	return total, nil
}

func (r *sqlxUserRepository) GetDailyWordCountLastNDays(userID, days int) (map[string]int, error) {
	result := make(map[string]int)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for i := 0; i < days; i++ {
		startDate := today.AddDate(0, 0, -i)
		endDate := startDate.AddDate(0, 0, 1)

		var dailyCount int
		query := `
			SELECT COALESCE(SUM(word_count), 0)
			FROM reading_records
			WHERE user_id = $1 AND read_at >= $2 AND read_at < $3
		`
		err := r.DB.Get(&dailyCount, query, userID, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get daily word count for %s: %w", startDate.Format("2006-01-02"), err)
		}

		result[startDate.Format("2006-01-02")] = dailyCount
	}
	return result, nil
}