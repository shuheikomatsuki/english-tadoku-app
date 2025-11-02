package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository" // ErrEmailAlreadyExists の比較用
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"    // service パッケージをインポート
)

func setupAuthTestHandler(t *testing.T) (*MockAuthService, *MockUserService, *echo.Echo) {
	mockAuthSvc := new(MockAuthService)
	mockUserSvc := new(MockUserService)
	e := echo.New()
	e.Validator = NewValidator()
	return mockAuthSvc, mockUserSvc, e
}

func TestAuthHandler_SignUp(t *testing.T) {
	mockAuthSvc, _, e := setupAuthTestHandler(t)
	h := NewAuthHandler(mockAuthSvc, nil)

	t.Run("success: should create a new user", func(t *testing.T) {
		email := fmt.Sprintf("signup-test-%d@example.com", time.Now().UnixNano())
		password := "password123"
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
		mockAuthSvc.On("SignUp", email, password).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code, "status code should be 201 Created")

		mockAuthSvc.AssertExpectations(t)
	})

	t.Run("fail: should return 409 Conflict if email exists", func(t *testing.T) {
		email := "exists@example.com"
		password := "password123"
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

		mockAuthSvc.On("SignUp", email, password).Return(repository.ErrEmailAlreadyExists).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		mockAuthSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockAuthSvc, _, e := setupAuthTestHandler(t)
	h := NewAuthHandler(mockAuthSvc, nil)

	email := "login-test@example.com"
	password := "password123"

	t.Run("success: should return a JWT token", func(t *testing.T) {
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
		expectedToken := "mocked.jwt.token"

		mockAuthSvc.On("ValidateUser", email, password).Return(testUser, nil).Once()
		mockAuthSvc.On("GenerateToken", testUser.ID).Return(expectedToken, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code, "status code should be 200 OK")

		var responseBody map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Contains(t, responseBody, "token")
		assert.Equal(t, expectedToken, responseBody["token"])

		mockAuthSvc.AssertExpectations(t)
	})

	t.Run("fail: should return 401 Unauthorized for invalid password", func(t *testing.T) {
		wrongPassword := "wrong-password"
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, wrongPassword)

		mockAuthSvc.On("ValidateUser", email, wrongPassword).Return(nil, fmt.Errorf("invalid password")).Once()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockAuthSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_GetUserStats(t *testing.T) {
	mockAuthSvc, mockUserSvc, e := setupAuthTestHandler(t)

	h := NewAuthHandler(mockAuthSvc, mockUserSvc)

	claims := &JwtCustomClaims{
		testUserID,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run("success: should return user stats", func(t *testing.T) {
		// Service 層が返すモデルを作成
		expectedStats := &service.UserStats{
			TotalWordCount:     1000,
			TodayWordCount:     100,
			WeeklyWordCount:    300,
			MonthlyWordCount:   500,
			YearlyWordCount:    800,
			Last7DaysWordCount: map[string]int{time.Now().Format("2006-01-02"): 100},
		}

		mockUserSvc.On("GetUserStats", testUserID).Return(expectedStats, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token) 

		err := h.GetUserStats(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response UserStatsResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.Equal(t, expectedStats.TotalWordCount, response.TotalWordCount)
		assert.Equal(t, expectedStats.TodayWordCount, response.TodayWordCount)
		assert.Equal(t, expectedStats.Last7DaysWordCount[time.Now().Format("2006-01-02")], response.Last7DaysWordCount[time.Now().Format("2006-01-02")])

		mockUserSvc.AssertExpectations(t)
	})
}
