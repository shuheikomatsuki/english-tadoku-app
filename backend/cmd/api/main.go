package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shuheikomatsuki/readoku/backend/internal/handler"
	authMiddleware "github.com/shuheikomatsuki/readoku/backend/internal/middleware"
	"github.com/shuheikomatsuki/readoku/backend/internal/repository"
	"github.com/shuheikomatsuki/readoku/backend/internal/service"
)

func main() {
	e := buildServer()

	// Lambda 環境では API Gateway (HTTP API) と接続するハンドラで起動
	if isLambda() {
		adapter := echoadapter.NewV2(e)
		lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return adapter.ProxyWithContext(ctx, req)
		})
		return
	}

	// ローカル / 常時稼働サーバーモード
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func buildServer() *echo.Echo {
	e := echo.New()

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		// ローカル開発用
		frontendURL = "http://localhost:5173"
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			frontendURL,
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))

	e.Validator = handler.NewValidator()

	// DB接続
	db, err := repository.NewDBConnection()
	if err != nil {
		e.Logger.Fatal("Failed to connect to database:", err)
	}

	limitStr := os.Getenv("DAILY_GENERATION_LIMIT")
	dailyLimit, err := strconv.Atoi(limitStr)
	if err != nil || dailyLimit <= 0 {
		log.Println("Invalid or missing DAILY_GENERATION_LIMIT, using default value (10)")
		dailyLimit = 10 // デフォルト値
	}
	e.Logger.Infof("Daily story generation limit set to %d", dailyLimit)

	// --- 依存関係の注入 ---

	// Repository層
	userRepo := repository.NewUserRepository(db)
	storyRepo := repository.NewStoryRepository(db)
	readingRecordRepo := repository.NewReadingRecordRepository(db)

	// Service層
	llmService, err := service.NewLLMService(os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		e.Logger.Fatal("Failed to init LLMService:", err)
	}
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(readingRecordRepo, userRepo, dailyLimit)
	storyService := service.NewStoryService(storyRepo, readingRecordRepo, userRepo, llmService, dailyLimit)

	// Handler層
	authHandler := handler.NewAuthHandler(authService, userService)
	storyHandler := handler.NewStoryHandler(storyService)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.HEAD("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/db-health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.Logger().Warnf("db-health check failed: %v", err)
			return c.String(http.StatusServiceUnavailable, "database connection unhealthy")
		}
		return c.String(http.StatusOK, "database connection healthy")
	})
	e.HEAD("/db-health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.Logger().Warnf("db-health check failed: %v", err)
			return c.String(http.StatusServiceUnavailable, "database connection unhealthy")
		}
		return c.String(http.StatusOK, "database connection healthy")
	})

	// ルーティング設定
	api := e.Group("/api/v1")
	api.POST("/signup", authHandler.SignUp)
	api.POST("/login", authHandler.Login)

	userRoutes := api.Group("/users")
	userRoutes.Use(authMiddleware.JWTAuthMiddleware)
	userRoutes.GET("/me/stats", authHandler.GetUserStats)
	userRoutes.GET("/me/generation-status", authHandler.GetGenerationStatus)

	// 認証が必要なグループ
	stories := api.Group("/stories")
	stories.Use(authMiddleware.JWTAuthMiddleware)
	stories.POST("", storyHandler.GenerateStory)
	stories.GET("", storyHandler.GetStories)
	stories.GET("/:id", storyHandler.GetStory)
	stories.DELETE("/:id", storyHandler.DeleteStory)
	stories.PATCH("/:id", storyHandler.UpdateStory)
	stories.POST("/:id/read", storyHandler.MarkStoryAsRead)
	stories.DELETE("/:id/read/latest", storyHandler.UndoLastRead)

	return e
}

func isLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
