package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// --- ヘルパー関数 ---

// テーブルをクリーンアップするヘルパー
func cleanupUsersTable(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("DELETE FROM users")
	require.NoError(t, err)
}

// --- テストケース ---

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)

	t.Run("CreateUser and FindUserByEmail", func(t *testing.T) {
		// 1. 準備
		password := "password123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		userToCreate := &model.User{
			Email:        fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
			PasswordHash: string(hashedPassword),
		}

		t.Cleanup(func() {
			cleanupUsersTable(t, db)
		})

		// 2. 実行
		// ユーザーを作成
		err = userRepo.CreateUser(userToCreate)
		require.NoError(t, err)
		assert.NotZero(t, userToCreate.ID)
		assert.NotZero(t, userToCreate.CreatedAt)

		// 作成したユーザーをEmailで検索
		foundUser, err := userRepo.FindUserByEmail(userToCreate.Email)

		// 3. 検証
		require.NoError(t, err)
		require.NotNil(t, foundUser)
		assert.Equal(t, userToCreate.ID, foundUser.ID)
		assert.Equal(t, userToCreate.Email, foundUser.Email)
		assert.Equal(t, userToCreate.PasswordHash, foundUser.PasswordHash)
	})

	t.Run("FindUserByEmail for non-existent user", func(t *testing.T) {
		// 1. 準備
		nonExistentEmail := "nonexistent@example.com"

		// 2. 実行
		foundUser, err := userRepo.FindUserByEmail(nonExistentEmail)

		// 3. 検証
		assert.Error(t, err, "should return an error for a non-existent user")
		assert.Nil(t, foundUser, "should return nil for a non-existent user")
	})
}