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