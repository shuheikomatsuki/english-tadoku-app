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
	GetUserStats(e echo.Context) error
}

type AuthHandler struct {
	UserRepo repository.IUserRepository
}

func NewAuthHandler(userRepo repository.IUserRepository) IAuthHandler {
	return &AuthHandler{
		UserRepo: userRepo,
	}
}

type UserStatsResponse struct {
	TotalWordCount     int            `json:"total_word_count"`
	TodayWordCount     int            `json:"today_word_count"`
	WeeklyWordCount    int            `json:"weekly_word_count"`
	MonthlyWordCount   int            `json:"monthly_word_count"`
	YearlyWordCount    int            `json:"yearly_word_count"`
	Last7DaysWordCount map[string]int `json:"last_7_days_word_count"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"`
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
		Email:        req.Email,
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

func (h *AuthHandler) GetUserStats(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	// 各種統計情報を計算
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayCount, err := h.UserRepo.GetWordCountInDateRange(userID, startOfDay, now.Add(24*time.Hour))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	// 週の始まりを月曜日とする
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜日の場合は7に調整
	}
	startOfWeek := startOfDay.AddDate(0, 0, -weekday+1)
	weeklyCount, err := h.UserRepo.GetWordCountInDateRange(userID, startOfWeek, startOfWeek.AddDate(0, 0, 7))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthlyCount, err := h.UserRepo.GetWordCountInDateRange(userID, startOfMonth, startOfMonth.AddDate(0, 1, 0))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	yearlyCount, err := h.UserRepo.GetWordCountInDateRange(userID, startOfYear, startOfYear.AddDate(1, 0, 0))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	totalWordCount, err := h.UserRepo.GetUserTotalWordCount(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	last7DaysCount, err := h.UserRepo.GetDailyWordCountLastNDays(userID, 7)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	if last7DaysCount == nil {
		last7DaysCount = make(map[string]int)
	}

	// レスポンス生成

	res := UserStatsResponse{
		TotalWordCount:     totalWordCount,
		TodayWordCount:     todayCount,
		WeeklyWordCount:    weeklyCount,
		MonthlyWordCount:   monthlyCount,
		YearlyWordCount:    yearlyCount,
		Last7DaysWordCount: last7DaysCount,
	}

	return c.JSON(http.StatusOK, res)
}
