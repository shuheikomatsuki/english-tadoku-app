package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	

	// .env から読み込むなどの初期設定

	// DB接続
	// db := repository.NewDB()

	// 依存関係を注入（DB）
	// userRepo := repository.NewUserRepository(db)
	// storyRepo := repository.NewStoryRepository(db)
	// llmService := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
	// authHandler := handler.NewAuthHandler(userRepo)
	// storyHandler := handler.NewStoryHandler(storyRepo, llmService)

	// ルーティング設定
	// api := e.Group("/api/v1")
	// api.POST("/signup", authHandler.SignUp)
	// api.POST("/login", authHandler.Login)
	

	// 認証が必要なグループ
	// stories := api.Group("/stories")
	// TODO: JWT ミドルウェアを追加する処理
	// stories.POST("", storyHandler.GenerateStory)
	// stories.GET("", storyHandler.GetStories)
	// stories.GET("/:id", storyHandler.GetStory))

	e.Logger.Fatal(e.Start(":8080"))
}