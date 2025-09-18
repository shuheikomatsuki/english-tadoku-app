package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// dsn := "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
	    dsn = "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err, "failed to connect to test database")
	require.NoError(t, db.Ping(), "failed to ping test database")
	return db
}

func createTestUser(t *testing.T, db *sqlx.DB) *model.User {
	uniqueEmail := fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())
	user := &model.User{
		Email:        uniqueEmail,
		PasswordHash: "hashed_password",
	}
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	err := db.QueryRowx(query, user.Email, user.PasswordHash).Scan(&user.ID)
	require.NoError(t, err, "failed to create test user")

	t.Cleanup(func() {
		// storiesテーブルとの依存関係のため、先にstoriesを削除
		_, _ = db.Exec("DELETE FROM stories WHERE user_id = $1", user.ID)
		_, err := db.Exec("DELETE FROM users WHERE id = $1", user.ID)
		require.NoError(t, err, "failed to delete test user after test")
	})
	return user
}