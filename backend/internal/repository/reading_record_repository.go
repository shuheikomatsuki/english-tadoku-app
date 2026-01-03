package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuheikomatsuki/readoku/backend/internal/model"
)

// IReadingRecordRepository: reading_records テーブルの操作インターフェース
type IReadingRecordRepository interface {
	CreateReadingRecord(userID, storyID, wordCount int) error
	CountReadingRecords(userID, storyID int) (int, error)
	GetLatestReadingRecord(userID, storyID int) (*model.ReadingRecord, error)
	DeleteReadingRecord(recordID int, userID int) error

	GetUserTotalWordCount(userID int) (int, error)
	GetWordCountInDateRange(userID int, start, end time.Time) (int, error)
	GetDailyWordCountLastNDays(userID, days int, anchorTime time.Time) (map[string]int, error)
}

type sqlxReadingRecordRepository struct {
	DB *sqlx.DB
}

func NewReadingRecordRepository(db *sqlx.DB) IReadingRecordRepository {
	return &sqlxReadingRecordRepository{DB: db}
}

func (r *sqlxReadingRecordRepository) CreateReadingRecord(userID, storyID, wordCount int) error {
	query := `
		INSERT INTO reading_records(user_id, story_id, word_count)
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, userID, storyID, wordCount)
	if err != nil {
		return fmt.Errorf("failed to create reading record: %w", err)
	}
	return nil
}

func (r *sqlxReadingRecordRepository) CountReadingRecords(userID, storyID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM reading_records WHERE user_id = $1 AND story_id = $2`
	err := r.DB.Get(&count, query, userID, storyID)
	if err != nil {
		return 0, fmt.Errorf("failed to count reading records: %w", err)
	}
	return count, nil
}

func (r *sqlxReadingRecordRepository) GetLatestReadingRecord(userID, storyID int) (*model.ReadingRecord, error) {
	var record model.ReadingRecord
	query := `
		SELECT * FROM reading_records
		WHERE user_id = $1 AND story_id = $2
		ORDER BY read_at DESC
		LIMIT 1
	`
	err := r.DB.Get(&record, query, userID, storyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest reading record: %w", err)
	}
	return &record, nil
}

func (r *sqlxReadingRecordRepository) DeleteReadingRecord(recordID int, userID int) error {
	query := `DELETE FROM reading_records WHERE id = $1 AND user_id = $2`
	_, err := r.DB.Exec(query, recordID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reading record: %w", err)
	}
	return nil
}

// --- UserRepository から移管されたメソッド ---

func (r *sqlxReadingRecordRepository) GetUserTotalWordCount(userID int) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(word_count), 0) FROM reading_records WHERE user_id = $1`

	err := r.DB.Get(&total, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user total word count: %w", err)
	}

	return total, nil
}

func (r *sqlxReadingRecordRepository) GetWordCountInDateRange(userID int, start, end time.Time) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(word_count), 0)
		FROM reading_records
		WHERE user_id = $1 AND read_at >= $2 AND read_at < $3
	`
	err := r.DB.Get(&total, query, userID, start, end)
	if err != nil {
		return 0, fmt.Errorf("failed to get word count in date range: %w", err)
	}
	return total, nil
}

func (r *sqlxReadingRecordRepository) GetDailyWordCountLastNDays(userID, days int, anchorTime time.Time) (map[string]int, error) {
	result := make(map[string]int)

	now := anchorTime
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Pre-fill the map so missing days are returned as zero.
	for i := 0; i < days; i++ {
		day := today.AddDate(0, 0, -i)
		result[day.Format("2006-01-02")] = 0
	}

	startDate := today.AddDate(0, 0, -(days - 1))
	endDate := today.AddDate(0, 0, 1)

	type dailyRow struct {
		Day   time.Time `db:"day"`
		Total int       `db:"total"`
	}

	query := `
		SELECT
			date_trunc('day', read_at) AS day,
			SUM(word_count) AS total
		FROM reading_records
		WHERE user_id = $1
		  AND read_at >= $2
		  AND read_at < $3
		GROUP BY day
	`

	var rows []dailyRow
	if err := r.DB.Select(&rows, query, userID, startDate, endDate); err != nil {
		return nil, fmt.Errorf("failed to get daily word count: %w", err)
	}

	for _, row := range rows {
		result[row.Day.Format("2006-01-02")] = row.Total
	}

	return result, nil
}
