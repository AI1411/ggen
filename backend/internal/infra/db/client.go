package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	applogger "g_gen/internal/infra/logger"
)

// Client インターフェース
type Client interface {
	Conn(ctx context.Context) *gorm.DB
	Close() error
	Ping(ctx context.Context) error
	Transaction(ctx context.Context, fn func(tx Client) error) error
}

// SQLLogger インターフェース
type SQLLogger interface {
	logger.Interface
}

// SQLHandler はDatabaseHandlerの実装
type SQLHandler struct {
	Driver *gorm.DB
}

// JSONLogger is a custom GORM logger that uses our application's JSON logger
type JSONLogger struct {
	logger        *applogger.Logger
	slowThreshold time.Duration
	logLevel      logger.LogLevel
}

// NewJSONLogger creates a new JSONLogger
func NewJSONLogger(appLogger *applogger.Logger) SQLLogger {
	// 環境変数でSQLログレベルを設定可能にする
	sqlLogLevel := logger.Info

	if logLevel := os.Getenv("SQL_LOG_LEVEL"); logLevel != "" {
		switch logLevel {
		case "debug":
			sqlLogLevel = logger.Info
		case "info":
			sqlLogLevel = logger.Info
		case "warn":
			sqlLogLevel = logger.Warn
		case "error":
			sqlLogLevel = logger.Error
		case "silent":
			sqlLogLevel = logger.Silent
		}
	}

	// スロークエリの閾値も環境変数で設定可能にする
	slowThreshold := time.Second

	if threshold := os.Getenv("SQL_SLOW_THRESHOLD"); threshold != "" {
		if duration, err := time.ParseDuration(threshold); err == nil {
			slowThreshold = duration
		}
	}

	return &JSONLogger{
		logger:        appLogger,
		slowThreshold: slowThreshold,
		logLevel:      sqlLogLevel,
	}
}

// LogMode sets the log level
func (l *JSONLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level

	return &newLogger
}

// Info logs info messages
func (l *JSONLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.logger.InfoContext(ctx, msg, "data", data)
	}
}

// Warn logs warn messages
func (l *JSONLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.logger.WarnContext(ctx, msg, "data", data)
	}
}

// Error logs error messages
func (l *JSONLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.logger.ErrorContext(ctx, errors.New(msg), msg, "data", data)
	}
}

// Trace logs SQL statements with execution time
func (l *JSONLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.logger.ErrorContext(ctx, err, "SQL error",
			"sql", sql,
			"rows", rows,
			"elapsed_ms", elapsed.Milliseconds())

		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.logger.WarnContext(ctx, "SLOW SQL",
			"sql", sql,
			"rows", rows,
			"elapsed_ms", elapsed.Milliseconds(),
			"threshold_ms", l.slowThreshold.Milliseconds())

		return
	}

	if l.logLevel >= logger.Info {
		l.logger.InfoContext(ctx, "SQL",
			"sql", sql,
			"rows", rows,
			"elapsed_ms", elapsed.Milliseconds())
	}
}

// DatabaseConfig データベース設定
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	Timezone        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// DefaultDatabaseConfig デフォルト設定を返す
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "postgres",
		Port:            "5432",
		User:            "postgres",
		Password:        "postgres",
		DBName:          "gen",
		SSLMode:         "disable",
		Timezone:        "Asia/Tokyo",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
	}
}

// NewSQLHandler creates a new SQLHandler
func NewSQLHandler(config *DatabaseConfig, appLogger *applogger.Logger) (Client, error) {
	if config == nil {
		config = DefaultDatabaseConfig()
	}

	// PostgreSQL形式の接続文字列
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode, config.Timezone)

	// PostgreSQLに直接接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewJSONLogger(appLogger),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続プールの設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 接続プールの設定
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &SQLHandler{
		Driver: db,
	}, nil
}

// Conn returns the underlying GORM DB instance
func (s *SQLHandler) Conn(ctx context.Context) *gorm.DB {
	return s.Driver.WithContext(ctx)
}

// Close closes the database connection
func (s *SQLHandler) Close() error {
	sqlDB, err := s.Driver.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// Ping verifies a connection to the database is still alive
func (s *SQLHandler) Ping(ctx context.Context) error {
	sqlDB, err := s.Driver.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.PingContext(ctx)
}

// Transaction executes a function within a database transaction
func (s *SQLHandler) Transaction(ctx context.Context, fn func(tx Client) error) error {
	return s.Driver.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txHandler := &SQLHandler{Driver: tx}
		return fn(txHandler)
	})
}

// MockDatabaseHandler テスト用のモック実装
type MockDatabaseHandler struct {
	db           *gorm.DB
	shouldError  bool
	errorMessage string
}

// NewMockDatabaseHandler creates a new mock database handler
func NewMockDatabaseHandler() Client {
	return &MockDatabaseHandler{
		db: &gorm.DB{}, // ダミーのGORMインスタンス
	}
}

// Conn returns the mock GORM DB instance
func (m *MockDatabaseHandler) Conn(ctx context.Context) *gorm.DB {
	return m.db.WithContext(ctx)
}

// Close closes the mock connection
func (m *MockDatabaseHandler) Close() error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	return nil
}

// Ping verifies mock connection
func (m *MockDatabaseHandler) Ping(ctx context.Context) error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	return nil
}

// Transaction executes a function within a mock transaction
func (m *MockDatabaseHandler) Transaction(ctx context.Context, fn func(tx Client) error) error {
	if m.shouldError {
		return errors.New(m.errorMessage)
	}

	return fn(m)
}

// SetError sets the mock to return errors
func (m *MockDatabaseHandler) SetError(shouldError bool, message string) {
	m.shouldError = shouldError
	m.errorMessage = message
}
