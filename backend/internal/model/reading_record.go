package model

import (
	"time"
)

type ReadingRecord struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	StoryID   int       `json:"story_id" db:"story_id"`
	WordCount int       `json:"word_count" db:"word_count"`
	ReadAt    time.Time `json:"read_at" db:"read_at"`
}
