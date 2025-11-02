package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- テストケース ---

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)

	userRepo := NewUserRepository(db)

	t.Run("CreateUser and FindUserByEmail", func(t *testing.T) {
		userToCreate := createTestUser(t, db)

		foundUser, err := userRepo.FindUserByEmail(userToCreate.Email)

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		assert.Equal(t, userToCreate.ID, foundUser.ID)
		assert.Equal(t, userToCreate.Email, foundUser.Email)
		assert.Equal(t, userToCreate.PasswordHash, foundUser.PasswordHash)
	})

	t.Run("FindUserByEmail for non-existent user", func(t *testing.T) {
		nonExistentEmail := "nonexistent@example.com"
		foundUser, err := userRepo.FindUserByEmail(nonExistentEmail)
		assert.Error(t, err, "should return an error for a non-existent user")
		assert.Nil(t, foundUser, "should return nil for a non-existent user")
	})
}