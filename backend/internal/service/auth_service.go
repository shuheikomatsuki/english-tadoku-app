package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shuheikomatsuki/readoku/backend/internal/model"
	"github.com/shuheikomatsuki/readoku/backend/internal/repository"
	"github.com/shuheikomatsuki/readoku/backend/internal/timeutil"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUp(email, password string) error
	ValidateUser(email, password string) (*model.User, error)
	GenerateToken(userID int) (string, error)
}

type AuthService struct {
	UserRepo repository.IUserRepository
}

func NewAuthService(userRepo repository.IUserRepository) IAuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return repository.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *AuthService) ValidateUser(email, password string) (*model.User, error) {
	// ユーザー検索
	user, err := s.UserRepo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// パスワード検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// パスワード不一致
		return nil, fmt.Errorf("invalid email or password")
	}
	return user, nil
}

type JwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(userID int) (string, error) {
	claims := &JwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeutil.NowTokyo().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return t, nil
}
