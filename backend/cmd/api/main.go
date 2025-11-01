package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	
	authMiddleware "github.com/shuheikomatsuki/english-tadoku-app/backend/internal/middleware"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/handler"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/repository"
	"github.com/shuheikomatsuki/english-tadoku-app/backend/internal/service"
)

func main() {
	e := echo.New()
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.Validator = handler.NewValidator()

	// .env から読み込むなどの初期設定

	// DB接続
	db, err := repository.NewDBConnection()
	if err != nil {
		e.Logger.Fatal("Failed to connect to database:", err)
	}

	// 依存関係を注入（DB）
	userRepo := repository.NewUserRepository(db)
	storyRepo := repository.NewStoryRepository(db)
	llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
	authHandler := handler.NewAuthHandler(userRepo)
	storyHandler := handler.NewStoryHandler(storyRepo, llmService)

	// ルーティング設定
	api := e.Group("/api/v1")
	api.POST("/signup", authHandler.SignUp)
	api.POST("/login", authHandler.Login)

	userRoutes := api.Group("/users")
	userRoutes.Use(authMiddleware.JWTAuthMiddleware)
	userRoutes.GET("/me/stats", authHandler.GetUserStats)
	
	// 認証が必要なグループ
	stories := api.Group("/stories")
	stories.Use(authMiddleware.JWTAuthMiddleware) // TODO: JWT ミドルウェアを追加する処理
	stories.POST("", storyHandler.GenerateStory)
	stories.GET("", storyHandler.GetStories)
	stories.GET("/:id", storyHandler.GetStory)
	stories.DELETE("/:id", storyHandler.DeleteStory)
	stories.PATCH("/:id", storyHandler.UpdateStory)
	stories.POST("/:id/read", storyHandler.MarkStoryAsRead)

	e.Logger.Fatal(e.Start(":8080"))
}