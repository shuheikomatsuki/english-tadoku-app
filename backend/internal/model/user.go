package model

import (
	"time"
)

type User struct {
	ID               int        `json:"id"            db:"id"`
	Email            string     `json:"email"         db:"email"`
	PasswordHash     string     `json:"password_hash" db:"password_hash"`
	GenerationCount  int        `json:"generation_count,omitempty"  db:"generation_count"`
	LastGenerationAt *time.Time `json:"last_generation_at,omitempty" db:"last_generation_at"`
	CreatedAt        time.Time  `json:"created_at"    db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"    db:"updated_at"`
}
