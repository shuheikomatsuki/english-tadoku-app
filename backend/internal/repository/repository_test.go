package repository

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// dsn := "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping to test database: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM reading_records")
		require.NoError(t, err, "failed to cleanup reading_records table")

		_, err = db.Exec("DELETE FROM stories")
		require.NoError(t, err, "failed to cleanup stories table")

		_, err = db.Exec("DELETE FROM users")
		require.NoError(t, err, "failed to cleanup users table")

		db.Close()
	})

	return db
}

// --- 共通テストヘルパー ---

func createTestUser(t *testing.T, db *sqlx.DB) *model.User {
	userRepo := NewUserRepository(db)

	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &model.User{
		Email:        fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
		PasswordHash: string(hashedPassword),
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err, "failed to create test user for setup")

	return user
}

func createTestStory(t *testing.T, db *sqlx.DB, userID int, title string, wordCount int) *model.Story {
	storyRepo := NewStoryRepository(db)
	story := &model.Story{
		UserID:    userID,
		Title:     title,
		Content:   fmt.Sprintf("Content for %s", title),
		WordCount: wordCount,
	}
	err := storyRepo.CreateStory(story)
	require.NoError(t, err, "failed to create test story for setup")

	return story
}

func createTestReadingRecord(t *testing.T, db *sqlx.DB, userID, storyID, wordCount int, readAt time.Time) *model.ReadingRecord {
	record := &model.ReadingRecord{
		UserID: userID,
		StoryID: storyID,
		WordCount: wordCount,
		ReadAt: readAt,
	}

	query := `
		INSERT INTO reading_records (user_id, story_id, word_count, read_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := db.QueryRowx(query, record.UserID, record.StoryID, record.WordCount, record.ReadAt).Scan(&record.ID)
	require.NoError(t, err, "failed to create test reading record for setup")

	return record
}
