// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (
	"github.com/labstack/echo/v4"
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

func (h *AuthHandler) SignUp(e echo.Context) error {
	return nil
}

func (h *AuthHandler) Login(e echo.Context) error {
	return nil
}