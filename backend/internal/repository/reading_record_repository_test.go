package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadingRecordRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReadingRecordRepository(db)

	user1 := createTestUser(t, db)
	story1 := createTestStory(t, db, user1.ID, "Story 1", 100)
	story2 := createTestStory(t, db, user1.ID, "Story 2", 200)

	now := time.Now().Truncate(24 * time.Hour)
	todayRecord1 := createTestReadingRecord(t, db, user1.ID, story1.ID, 100, now.Add(1 * time.Hour))
	todayRecord2 := createTestReadingRecord(t, db, user1.ID, story2.ID, 200, now.Add(2 * time.Hour))

	yesterdayRecord := createTestReadingRecord(t, db, user1.ID, story1.ID, 100, now.Add(-10 * time.Hour))
	twoDaysAgoRecord := createTestReadingRecord(t, db, user1.ID, story1.ID, 100, now.Add(-54 * time.Hour))

	if yesterdayRecord == nil || twoDaysAgoRecord == nil {}

	user2 := createTestUser(t, db)
	story3 := createTestStory(t, db, user2.ID, "Story 3", 50)
	user2Record := createTestReadingRecord(t, db, user2.ID, story3.ID, 50, now)

	t.Run("CreateReadingRecord", func(t *testing.T) {

		user := createTestUser(t, db)
		story := createTestStory(t, db, user.ID, "Test Story for Create", 150)

		err := repo.CreateReadingRecord(user.ID, story.ID, story.WordCount)
		require.NoError(t, err)

		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM reading_records WHERE user_id = $1 AND story_id = $2", user.ID, story.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("CountReadingRecords", func(t *testing.T) {
		// user1, story1 は3件
		count, err := repo.CountReadingRecords(user1.ID, story1.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, count)

		// user1, story2 は1件
		count, err = repo.CountReadingRecords(user1.ID, story2.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// user2, story1 は0件
		count, err = repo.CountReadingRecords(user2.ID, story1.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("GetLatestReadingRecord", func(t *testing.T) {
		latest, err := repo.GetLatestReadingRecord(user1.ID, story1.ID)
		require.NoError(t, err)
		assert.Equal(t, todayRecord1.ID, latest.ID)
		assert.Equal(t, todayRecord1.WordCount, latest.WordCount)
	})

	t.Run("DeleteReadingRecord", func(t *testing.T) {
		err := repo.DeleteReadingRecord(todayRecord2.ID, user1.ID)
		require.NoError(t, err)

		// 削除されたか確認
		_, err = repo.GetLatestReadingRecord(todayRecord2.ID, user1.ID)
		assert.Error(t, err, "should be no rows after delete")

		// 他のユーザーの記録は削除できない
		err = repo.DeleteReadingRecord(user2Record.ID, user1.ID)
		require.NoError(t, err)

		count, err := repo.CountReadingRecords(user2.ID, story3.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should not delete other user's record")
	})

	t.Run("GetUserTotalWordCount", func(t *testing.T) {
		// user1: 100(today) + 200(today) + 100(yesterday) + 100(two days ago) = 500
		// (todayRecord2 は DeleteReadingRecord で削除されたので 300)
		total, err := repo.GetUserTotalWordCount(user1.ID)
		require.NoError(t, err)
		assert.Equal(t, 300, total) // 100 + 100 + 100

		// user2: 50
		total, err = repo.GetUserTotalWordCount(user2.ID)
		require.NoError(t, err)
		assert.Equal(t, 50, total)
	})

	t.Run("GetWordCountInDateRange", func(t *testing.T) {
		startOfDay := now
		endOfDay := now.Add(24 * time.Hour)
		total, err := repo.GetWordCountInDateRange(user1.ID, startOfDay, endOfDay)
		require.NoError(t, err)
		assert.Equal(t, 100, total)

		startOfYesterday := now.Add(-24 * time.Hour)
		endOfYesterday := now
		total, err = repo.GetWordCountInDateRange(user1.ID, startOfYesterday, endOfYesterday)
		require.NoError(t, err)
		assert.Equal(t, 100, total)
	})

	t.Run("GetDailyWordCountLastNDays", func(t *testing.T) {
		result, err := repo.GetDailyWordCountLastNDays(user1.ID, 3)
		require.NoError(t, err)

		todayKey := now.Format("2006-01-02")
		yesterdayKey := now.AddDate(0, 0, -1).Format("2006-01-02")
		twoDaysAgoKey := now.AddDate(0, 0, -2).Format("2006-01-02")

		assert.Equal(t, 100, result[todayKey], "today's count")
		assert.Equal(t, 100, result[yesterdayKey], "yesterday's count")
		assert.Equal(t, 100, result[twoDaysAgoKey], "two days ago's count")
	})
}