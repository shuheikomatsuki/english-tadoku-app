package handler

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/model"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
)

type IAuthHandler interface {
	SignUp(e echo.Context) error
	Login(e echo.Context) error
}

type AuthHandler struct {
	UserRepo repository.IUserRepository
}

func NewAuthHandler(userRepo repository.IUserRepository) IAuthHandler {
	return &AuthHandler{
		UserRepo: userRepo,
	}
}

type SignUpRequest struct {
	Email 	 string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email 	 string `json:"email"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to hash password")
	}

	user := &model.User{
		Email: req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.UserRepo.CreateUser(user); err != nil {
		// TODO: email が重複した際のエラーハンドリング
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "This email address is already registered."})
		}
		return c.JSON(http.StatusInternalServerError, "failed to create user")
	}

	return c.JSON(http.StatusCreated, "user created successfully")
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	user, err := h.UserRepo.FindUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}

	claims := &JwtCustomClaims{
		user.ID, 
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}