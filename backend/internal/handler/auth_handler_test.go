package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
)

func TestAuthHandler_SignUp(t *testing.T) {
	userRepo := repository.NewUserRepository(testDB)
	h := NewAuthHandler(userRepo)
	e := echo.New()

	cleanupUserTable(t)

	t.Run("success: should create a new user", func(t *testing.T) {
		email := fmt.Sprintf("signup-test-%d@example.com", time.Now().UnixNano())
		password := "password123"
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}:`, email, password)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SignUp(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code, "status code should be 201 Created")

		createdUser, err := userRepo.FindUserByEmail(email)
		require.NoError(t, err, "user should be found in the database")
		require.NotNil(t, createdUser)
		assert.Equal(t, email, createdUser.Email)
		assert.NotEqual(t, password, createdUser.PasswordHash)
	})

	// TODO: Emailが重複する場合など、失敗ケースのテストも追加する
}

func TestAuthHandler_Login(t *testing.T) {
	userRepo := repository.NewUserRepository(testDB)
	h := NewAuthHandler(userRepo)
	e := echo.New()

	cleanupUserTable(t)

	// テスト用のユーザーを事前に準備
	email := fmt.Sprintf("login-test-%d@example.com", time.Now().UnixNano())
	password := "password123"

	signUpReqBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", strings.NewReader(signUpReqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.SignUp(c))

	t.Run("success: should return a JWT token", func(t *testing.T) {
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)
		
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
		assert.NotEmpty(t, responseBody["token"])
	})

	t.Run("fail: should return 401 Unauthorized for invalid password", func(t *testing.T) {
		wrongPassword := "wrong-password"
		requestBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, wrongPassword)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}