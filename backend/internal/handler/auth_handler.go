package handler

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/shuheikomatsuki/readoku/backend/internal/repository"
	"github.com/shuheikomatsuki/readoku/backend/internal/service"
)

type IAuthHandler interface {
	SignUp(e echo.Context) error
	Login(e echo.Context) error
	GetUserStats(e echo.Context) error
	GetGenerationStatus(e echo.Context) error
}

type AuthHandler struct {
	AuthService service.IAuthService
	UserService service.IUserService
}

func NewAuthHandler(authSvc service.IAuthService, userSvc service.IUserService) IAuthHandler {
	return &AuthHandler{
		AuthService: authSvc,
		UserService: userSvc,
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

func getUserIDFromContext(c echo.Context) (int, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("failed to get user from context")
	}

	claims, ok := user.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, errors.New("failed to get claims from token")
	}

	return claims.UserID, nil
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	err := h.AuthService.SignUp(req.Email, req.Password)
	if err != nil {
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

	// (簡易バリデーション)
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}

	// ユーザー検証
	user, err := h.AuthService.ValidateUser(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}

	// トークン生成
	t, err := h.AuthService.GenerateToken(user.ID)
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

	stats, err := h.UserService.GetUserStats(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user stats"})
	}

	// レスポンス生成
	res := UserStatsResponse{
		TotalWordCount:     stats.TotalWordCount,
		TodayWordCount:     stats.TodayWordCount,
		WeeklyWordCount:    stats.WeeklyWordCount,
		MonthlyWordCount:   stats.MonthlyWordCount,
		YearlyWordCount:    stats.YearlyWordCount,
		Last7DaysWordCount: stats.Last7DaysWordCount,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetGenerationStatus(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid token")
	}

	status, err := h.UserService.GetGenerationStatus(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get generation status"})
	}

	return c.JSON(http.StatusOK, status)
}
