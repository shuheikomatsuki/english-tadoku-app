package service

import (
	"testing"

	"github.com/shuheikomatsuki/readoku/backend/internal/model"
	"github.com/shuheikomatsuki/readoku/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_SignUp(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	authService := NewAuthService(mockUserRepo)

	t.Run("success: should hash password and create user", func(t *testing.T) {
		email := "newuser@example.com"
		password := "password123"

		mockUserRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*model.User)
			assert.Equal(t, email, userArg.Email)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(userArg.PasswordHash), []byte(password)))
		}).Return(nil).Once()

		err := authService.SignUp(email, password)

		require.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail: should return error if email already exists", func(t *testing.T) {
		mockUserRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(repository.ErrEmailAlreadyExists).Once()

		err := authService.SignUp("exists@example.com", "password123")

		require.Error(t, err)
		assert.Equal(t, repository.ErrEmailAlreadyExists, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	authService := NewAuthService(mockUserRepo)

	t.Run("success: should return user if password matches", func(t *testing.T) {
		// testUser は service_test.go 内で定義した
		mockUserRepo.On("FindUserByEmail", testUser.Email).Return(testUser, nil).Once()

		user, err := authService.ValidateUser(testUser.Email, "password123")

		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail: should return error if password mismatches", func(t *testing.T) {
		mockUserRepo.On("FindUserByEmail", testUser.Email).Return(testUser, nil).Once()

		user, err := authService.ValidateUser(testUser.Email, "wrongpassword")

		require.Error(t, err)
		assert.Nil(t, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	authService := NewAuthService(mockUserRepo)

	// os.Setenv("JWT_SECRET", "test_secret_key_for_auth_service")

	t.Run("success: should generate a valid token", func(t *testing.T) {
		tokenString, err := authService.GenerateToken(testUser.ID)

		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)
	})
}
