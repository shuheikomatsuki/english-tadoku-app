package repository

import (
	"github.com/jmoiron/sqlx"

	// "github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type IUserRepository interface {
	// GetUserByID(id int) (*model.User, error)
	// GetUserByEmail(email string) (*model.User, error)
	// CreateUser(user *model.User) (int, error) // 作成したユーザーのIDを返す
	// UpdateUserPassword(id int, passwordHash string) error
	// DeleteUser(id int) error
}

type sqlxUserRepository struct {
	DB *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) IUserRepository {
	return &sqlxUserRepository{DB: db}
}