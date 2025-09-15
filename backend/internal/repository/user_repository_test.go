package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
)

// テスト用のDB接続をセットアップする
func setupTestDB(t *testing.T) *sqlx.DB {
	dsn := "host=localhost port=5432 user=postgres password=password dbname=tadoku_db sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err, "failed to connect to test database")

	err = db.Ping()
	require.NoError(t, err, "failed to ping test database")

	return db
}

func newTestStory(userID int) *model.Story {
	return &model.Story{
		UserID: userID,
		Title: "Test Story",
		Content: "This is a test story content.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ヘルパー：テスト用のユーザーを作成
func createTestUser(t *testing.T, db *sqlx.DB) *model.User {
	uniqueEmail := fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())

	user := &model.User{
		Email: uniqueEmail,
		PasswordHash: "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := db.QueryRowx(query, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	require.NoError(t, err, "failed to create test user")

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users WHERE id = $1", user.ID)
		require.NoError(t, err, "failed to delete test user")
	})

	return user
}

// ヘルパー：テスト用のストーリーを作成し、DBに保存する（テストの準備用）
func createTestStory(t *testing.T, db *sqlx.DB, storyRepo IStoryRepository, userID int, title string) *model.Story {
	story := &model.Story{
		UserID:  userID,
		Title:   title,
		Content: fmt.Sprintf("Content for %s", title),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := storyRepo.CreateStory(story)
	require.NoError(t, err, "failed to create test story for setup")

	// このテストが終了したときに、作成したストーリーを削除するように予約する
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM stories WHERE id = $1", story.ID)
		require.NoError(t, err, "failed to delete test story after test")
	})

	return story
}

func TestStoryRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	storyRepo := NewStoryRepository(db)
	// userRepo := NewUserRepository(db)

	t.Run("CreateStory", func(t *testing.T) {
		user := createTestUser(t, db)
		storyToCreate := newTestStory(user.ID)

		err := storyRepo.CreateStory(storyToCreate)

		require.NoError(t, err, "failed to create storyToCreate")
		assert.NotZero(t, storyToCreate.ID, "story ID should be set after creation")

		var fetchedStory model.Story
		err = db.Get(&fetchedStory, "SELECT * FROM stories WHERE id = $1", storyToCreate.ID)

		require.NoError(t, err, "failed to fetch created story from DB")
		assert.Equal(t, user.ID, fetchedStory.UserID, "story UserID should match test user ID")
	})

	t.Run("GetUserStories", func(t *testing.T) {
		// user := createTestUser(t, db)
		// storiesToGet := newTestStory(user.ID)
		// err := storyRepo.CreateStory(storiesToGet)
		// require.NoError(t, err, "failed to create story for GetStories test")

		// fetchedStories, err := storyRepo.GetUserStories(user.ID, 10, 0)

		// require.NoError(t, err, "failed to get user stories")
		// require.Len(t, fetchedStories, 1, "should return one story")
		// assert.Equal(t, storiesToGet.Title, fetchedStories[0].Title, "story title should match")

		userA := createTestUser(t, db)
		userB := createTestUser(t, db)

		storyA1 := createTestStory(t, db, storyRepo, userA.ID, "User A - Story 1")
		time.Sleep(1 * time.Millisecond)
		storyA2 := createTestStory(t, db, storyRepo, userA.ID, "User A - Story 2")
		time.Sleep(1 * time.Millisecond)
		storyA3 := createTestStory(t, db, storyRepo, userA.ID, "User A - Story 3")
		time.Sleep(1 * time.Millisecond)
		storyB1 := createTestStory(t, db, storyRepo, userB.ID, "User B - Story 1")

		fetchedStoriesA, err := storyRepo.GetUserStories(userA.ID, 10, 0)
		require.NoError(t, err, "failed to get user A stories")
		require.Len(t, fetchedStoriesA, 3, "should return three stories for user A")

		for _, story := range fetchedStoriesA {
			assert.Equal(t, userA.ID, story.UserID, "story UserID should match User A ID")
		}

		assert.Equal(t, storyA3.Title, fetchedStoriesA[0].Title, "most recent story should be first")
		assert.Equal(t, storyA2.Title, fetchedStoriesA[1].Title, "second most recent story should be second")
		assert.Equal(t, storyA1.Title, fetchedStoriesA[2].Title, "oldest story should be last")

		fetchedStoriesB, err := storyRepo.GetUserStories(userB.ID, 10, 0)
		require.NoError(t, err, "failed to get user B stories")
		require.Len(t, fetchedStoriesB, 1, "should return one story for user B")
		assert.Equal(t, userB.ID, fetchedStoriesB[0].UserID, "story UserID should match User B ID")
		assert.Equal(t, storyB1.Title, fetchedStoriesB[0].Title, "story title should match for user B")
	})

	t.Run("GetUserStory", func(t *testing.T) {
		user := createTestUser(t, db)
		storyToGet := newTestStory(user.ID)
		err := storyRepo.CreateStory(storyToGet)
		require.NoError(t, err, "failed to create story for GetStory test")

		fetchedStory, err := storyRepo.GetUserStory(storyToGet.ID, user.ID)
		require.NoError(t, err, "failed to get user story")
		require.NotNil(t, fetchedStory, "story should not be nil")
		assert.Equal(t, storyToGet.ID, fetchedStory.ID, "story ID should match")
		assert.Equal(t, storyToGet.UserID, fetchedStory.UserID, "story UserID should match")
		assert.Equal(t, storyToGet.Content, fetchedStory.Content, "story content should match")
		assert.Equal(t, user.ID, fetchedStory.UserID, "story UserID should match test user ID")
	})

	t.Run("DeleteStory", func(t *testing.T) {
		user := createTestUser(t, db)
		storyToDelete := newTestStory(user.ID)
		err := storyRepo.CreateStory(storyToDelete)
		require.NoError(t, err, "failed to create story for DeleteStory test")
		
		err = storyRepo.DeleteStory(storyToDelete.ID)
		require.NoError(t, err, "failed to delete story")

		// 削除されたことを確認
		deletedStory, err := storyRepo.GetUserStory(storyToDelete.ID, user.ID)
		assert.Error(t, err, "should return error for deleted story")
		assert.Nil(t, deletedStory, "should return nil for a deleted story")
	})
}