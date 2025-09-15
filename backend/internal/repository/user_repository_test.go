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

// ヘルパー：テスト用のストーリーを作成し、DBに保存する
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
		storyCounts := []int{3, 4, 2, 1} // 各ユーザーのストーリ数
		numUsers := len(storyCounts)
		const limit = 1
		const offset = 0

		users := make([]*model.User, numUsers)
		allStories := make([][]*model.Story, numUsers) // 各ユーザーのストーリを格納するスライス

		for i, numStories := range storyCounts {
			users[i] = createTestUser(t, db)
			allStories[i] = make([]*model.Story, numStories)

			for j := 0; j < numStories; j++ {
				title := fmt.Sprintf("User %d - Story %d", i+1, j+1)
				allStories[i][j] = createTestStory(t, db, storyRepo, users[i].ID, title)
				if j < numStories - 1 {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}

		for i, user := range users {
			t.Run(fmt.Sprintf("for User %d", i+1), func(t * testing.T) {
				expectedStories := allStories[i]
				expectedCount := storyCounts[i]

				// ユーザーのストーリを取得
				fetchedStories, err := storyRepo.GetUserStories(user.ID, limit, offset)
				require.NoError(t, err, "failed to get User %d stories", i+1)

				// 件数、所有者、順序を検証
				if expectedCount > limit {
					expectedCount = limit
				}
				assert.Len(t, fetchedStories, expectedCount, "should return %d stories for User %d", expectedCount, i+1)

				for j, fetchedStory := range fetchedStories {
					assert.Equal(t, user.ID, fetchedStory.UserID, "story UserID should match User %d ID", i+1)
					expectedStory := expectedStories[len(expectedStories) - 1 - j]
					assert.Equal(t, expectedStory.ID, fetchedStory.ID, "story IDs should match")
					assert.Equal(t, expectedStory.UserID, fetchedStory.UserID, "story UserIDs should match")
					assert.Equal(t, expectedStory.Title, fetchedStory.Title, "story titles should be in descending order by CreatedAt")
				}
			})
		}
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