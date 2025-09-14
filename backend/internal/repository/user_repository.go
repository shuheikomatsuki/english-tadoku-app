package repository

import (
	"github.com/jmoiron/sqlx"

	// "github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

type IUserRepository interface {
}

type IStoryRepository interface {
	// Create(story *model.Story) (string, error)
}

type sqlxUserRepository struct {
	DB *sqlx.DB
}

type sqlxStoryRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &sqlxUserRepository{DB: db}
}

func NewStoryRepository(db *sqlx.DB) IStoryRepository {
	return &sqlxStoryRepository{DB: db}
}