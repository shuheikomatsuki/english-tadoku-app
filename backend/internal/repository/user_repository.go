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
	GetUserByID(userID int) (*model.User, error)
	UpdateGenerationStatus(userID int, newCount int, newDate time.Time) error
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

func (r *sqlxUserRepository) GetUserByID(userID int) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.DB.Get(&user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id %w", err)
	}
	return &user, nil
}

func (r *sqlxUserRepository) UpdateGenerationStatus(userID, newCount int, newDate time.Time) error {
	query := `
		UPDATE users
		SET generation_count = $1, last_generation_at = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.DB.Exec(query, newCount, newDate, userID)
	if err != nil {
		return fmt.Errorf("failed to update generation status: %w", err)
	}
	return nil
}
