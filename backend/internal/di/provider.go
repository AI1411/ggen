package di

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	domain "g_gen/internal/domain/repository"
	"g_gen/internal/env"
	"g_gen/internal/handler"
	"g_gen/internal/infra/datastore"
	"g_gen/internal/infra/db"
	"g_gen/internal/infra/logger"
	"g_gen/internal/server/middleware"
	"g_gen/internal/usecase"
)

// ProvideLogger creates a new logger instance
func ProvideLogger() *logger.Logger {
	return logger.New(logger.DefaultConfig())
}

// ProvideEnvValues creates a new env values instance
func ProvideEnvValues() (*env.Values, error) {
	return env.NewValues()
}

// ProvideDBClient creates a new database client
func ProvideDBClient(lc fx.Lifecycle, l *logger.Logger) (db.Client, error) {
	e, err := env.NewValues()
	if err != nil {
		l.Error("failed to load environment variables", "error", err)
		return nil, err
	}
	dbClient, err := db.NewSQLHandler(&db.DatabaseConfig{
		Host:            e.DatabaseHost,
		Port:            e.DatabasePort,
		User:            e.DatabaseUsername,
		Password:        e.DatabasePassword,
		DBName:          e.DatabaseName,
		SSLMode:         "disable",
		Timezone:        "Asia/Tokyo",
		MaxIdleConns:    e.ConnectionMaxIdle,
		MaxOpenConns:    e.ConnectionMaxOpen,
		ConnMaxLifetime: e.ConnectionMaxLifetime,
	}, l)
	if err != nil {
		l.Error("failed to connect to database", "error", err)
		return nil, err
	}

	// Register lifecycle hooks for the database client
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			l.Info("Closing database connection")
			return nil // Add proper cleanup if needed
		},
	})

	return dbClient, nil
}

// ProvideGinEngine creates and configures a new Gin engine
func ProvideGinEngine(l *logger.Logger) *gin.Engine {
	r := gin.Default()

	// ミドルウェアの設定
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.NewLogging(l))
	r.Use(middleware.CORSMiddleware())

	return r
}

// ProvidePrefectureRepository creates a new prefecture repository
func ProvidePrefectureRepository(dbClient db.Client) domain.PrefectureRepository {
	ctx := context.Background()
	return datastore.NewPrefectureRepository(ctx, dbClient)
}

// ProvidePrefectureUseCase creates a new prefecture use case
func ProvidePrefectureUseCase(repo domain.PrefectureRepository) usecase.PrefectureUseCase {
	return usecase.NewPrefectureUseCase(repo)
}

// ProvidePrefectureHandler creates a new prefecture handler
func ProvidePrefectureHandler(l *logger.Logger, prefectureUseCase usecase.PrefectureUseCase) handler.PrefectureHandler {
	return handler.NewPrefectureHandler(l, prefectureUseCase)
}

func Provider() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideLogger,
			ProvideEnvValues,
			ProvideDBClient,
			ProvideGinEngine,
			ProvidePrefectureRepository,
			ProvidePrefectureUseCase,
			ProvidePrefectureHandler,
		),
	)
}
