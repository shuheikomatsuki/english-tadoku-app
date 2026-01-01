package repository

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shuheikomatsuki/readoku/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- テストケース ---

func TestStoryRepository(t *testing.T) {
	db := setupTestDB(t)

	storyRepo := NewStoryRepository(db)

	t.Run("CreateStory", func(t *testing.T) {
		const (
			title   = "A Brand New Story"
			content = "Content for the new story."
		)

		user := createTestUser(t, db)
		wordCount := len(strings.Fields(content))

		storyToCreate := &model.Story{
			UserID:    user.ID,
			Title:     title,
			Content:   content,
			WordCount: wordCount,
		}

		err := storyRepo.CreateStory(storyToCreate)

		require.NoError(t, err)
		assert.NotZero(t, storyToCreate.ID)

		var fetchedStory model.Story
		err = db.Get(&fetchedStory, "SELECT * FROM stories WHERE id = $1", storyToCreate.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, fetchedStory.UserID)
		assert.Equal(t, title, fetchedStory.Title)
		assert.Equal(t, content, fetchedStory.Content)
		assert.NotZero(t, fetchedStory.CreatedAt)
		assert.NotZero(t, fetchedStory.UpdatedAt)
		assert.Equal(t, storyToCreate.WordCount, fetchedStory.WordCount)
	})

	t.Run("GetUserStories", func(t *testing.T) {
		storyCounts := []int{3, 4, 2, 1}
		numUsers := len(storyCounts)

		const (
			limit  = 2
			offset = 0
		)

		users := make([]*model.User, numUsers)
		allStories := make([][]*model.Story, numUsers)

		for i, numStories := range storyCounts {
			users[i] = createTestUser(t, db)
			allStories[i] = make([]*model.Story, numStories)

			for j := 0; j < numStories; j++ {
				title := fmt.Sprintf("User %d - Story %d", i+1, j+1)
				allStories[i][j] = createTestStory(t, db, users[i].ID, title, 0)
				if j < numStories-1 {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}

		for i, user := range users {
			t.Run(fmt.Sprintf("for User %d", i+1), func(t *testing.T) {
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
					expectedStory := expectedStories[len(expectedStories)-1-j]
					assert.Equal(t, expectedStory.ID, fetchedStory.ID, "story IDs should match")
					assert.Equal(t, expectedStory.UserID, fetchedStory.UserID, "story UserIDs should match")
					assert.Equal(t, expectedStory.Title, fetchedStory.Title, "story titles should be in descending order by CreatedAt")
				}
			})
		}
	})

	t.Run("GetUserStory", func(t *testing.T) {
		user := createTestUser(t, db)
		storyToGet := createTestStory(t, db, user.ID, "A story to get", 10)

		fetchedStory, err := storyRepo.GetUserStory(storyToGet.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedStory)
		assert.Equal(t, storyToGet.ID, fetchedStory.ID)
		assert.Equal(t, storyToGet.Content, fetchedStory.Content)
	})

	t.Run("DeleteStory", func(t *testing.T) {
		user := createTestUser(t, db)
		storyToDelete := createTestStory(t, db, user.ID, "A story to delete", 10)

		err := storyRepo.DeleteStory(storyToDelete.ID)
		require.NoError(t, err)

		// 削除されたことを確認
		_, err = storyRepo.GetUserStory(storyToDelete.ID, user.ID)
		assert.Error(t, err)
	})

	t.Run("UpdateStoryTitle", func(t *testing.T) {
		user := createTestUser(t, db)
		originalStory := createTestStory(t, db, user.ID, "Original Title", 10)
		newTitle := "Updated Title"

		updatedStory, err := storyRepo.UpdateStoryTitle(originalStory.ID, user.ID, newTitle)

		require.NoError(t, err)
		assert.Equal(t, originalStory.ID, updatedStory.ID)
		assert.Equal(t, user.ID, updatedStory.UserID)
		assert.Equal(t, newTitle, updatedStory.Title)
		assert.Equal(t, originalStory.Content, updatedStory.Content)
		assert.Equal(t, originalStory.CreatedAt, updatedStory.CreatedAt)
		assert.True(t, updatedStory.UpdatedAt.After(originalStory.UpdatedAt), "UpdatedAt should be updated to a later time")

		var fetchedStory model.Story
		err = db.Get(&fetchedStory, "SELECT * FROM stories WHERE id = $1", originalStory.ID)
		require.NoError(t, err)
		assert.Equal(t, newTitle, fetchedStory.Title)
		assert.Equal(t, originalStory.Content, fetchedStory.Content)
		assert.Equal(t, originalStory.CreatedAt, fetchedStory.CreatedAt)
		assert.True(t, fetchedStory.UpdatedAt.After(originalStory.UpdatedAt), "UpdatedAt in DB should be updated to a later time")
	})

	// t.Run("CreateReadingRecord", func(t *testing.T) {
	// 	user := createTestUser(t, db)
	// 	story := createTestStory(t, storyRepo, user.ID, "This is a test story for reading record")

	// 	story.WordCount = len(strings.Fields(story.Content))

	// 	err := storyRepo.CreateReadingRecord(user.ID, story.ID, story.WordCount)

	// 	require.NoError(t, err)

	// 	var record model.ReadingRecord
	// 	err = db.Get(&record, "SELECT * FROM reading_records WHERE user_id = $1 AND story_id = $2", user.ID, story.ID)
	// 	require.NoError(t, err)
	// 	assert.Equal(t, user.ID, record.UserID)
	// 	assert.Equal(t, story.ID, record.StoryID)
	// 	assert.Equal(t, story.WordCount, record.WordCount)
	// 	assert.NotZero(t, record.ReadAt)
	// })
}
