package model

import (
	"time"
)

type Story struct {
	ID        int       `json:"id"         db:"id"`
	UserID    int       `json:"user_id"    db:"user_id"`
	Title     string    `json:"title"      db:"title"`
	Content   string    `json:"content"    db:"content"`
	WordCount int       `json:"word_count" db:"word_count"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
