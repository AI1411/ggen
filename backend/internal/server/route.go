package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"g_gen/internal/env"
	"g_gen/internal/handler"
	"g_gen/internal/infra/db"
	"g_gen/internal/infra/logger"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	lc fx.Lifecycle,
	r *gin.Engine,
	l *logger.Logger,
	dbClient db.Client,
	env *env.Values,
	prefectureHandler handler.PrefectureHandler,
) {
	// Context for health check
	ctx := context.Background()
	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		if err := dbClient.Ping(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "unhealthy",
				"error":  "database ping failed",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// 都道府県関連のルート
	r.GET("/prefectures", prefectureHandler.ListPrefectures)
	r.GET("/prefectures/:code", prefectureHandler.GetPrefecture)

	// Swagger JSON エンドポイント
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("./docs/api/swagger.json")
	})

	// Register lifecycle hooks for the HTTP server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				l.Info(fmt.Sprintf("Starting server on :%s", env.ServerPort))
				if err := r.Run(fmt.Sprintf(":%s", env.ServerPort)); err != nil {
					l.Error("Failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Shutting down server")
			return nil // Add proper cleanup if needed
		},
	})
}
