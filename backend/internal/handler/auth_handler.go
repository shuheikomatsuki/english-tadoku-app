// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
)

type IAuthHandler interface {
	SignUp() error
	Login() error
}

type AuthHandler struct {
	UserRepo repository.IUserRepository
}

func NewAuthHandler(userRepo repository.IUserRepository) IAuthHandler {
	return &AuthHandler{
		UserRepo: userRepo,
	}
}

func (h *AuthHandler) SignUp() error {
	return nil
}

func (h *AuthHandler) Login() error {
	return nil
}