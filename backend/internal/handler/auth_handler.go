// userRepo := repository.NewUserRepository(db)
// storyRepo := repository.NewStoryRepository(db)
// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
// authHandler := handler.NewAuthHandler(userRepo)
// こんな感じでmain.goで使う想定

package handler

import (

)

type IAuthHandler interface {
	// SignUp() error
	// Login() error
}

type AuthHandler struct {

}

func NewAuthHandler() IAuthHandler {
	return &AuthHandler{}
}