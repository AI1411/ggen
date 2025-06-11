package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

// LogLevel represents the logging level.
type LogLevel string

const (
	// DebugLevel logs debug messages.
	DebugLevel LogLevel = "debug"
	// InfoLevel logs informational messages.
	InfoLevel LogLevel = "info"
	// WarnLevel logs warning messages.
	WarnLevel LogLevel = "warn"
	// ErrorLevel logs error messages.
	ErrorLevel LogLevel = "error"
)

// Logger is a wrapper around slog.Logger that provides additional functionality.
type Logger struct {
	*slog.Logger
}

// Config holds the configuration for the logger.
type Config struct {
	// Level is the minimum log level that will be logged.
	Level LogLevel
	// Output is where the logs will be written to.
	Output io.Writer
	// AddSource adds the source file and line number to the log.
	AddSource bool
	// JSON determines whether the logs should be formatted as JSON.
	JSON bool
}

// DefaultConfig returns a default configuration for the logger.
func DefaultConfig() Config {
	// 環境変数でログレベルを設定可能にする
	level := InfoLevel

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch logLevel {
		case "debug":
			level = DebugLevel
		case "info":
			level = InfoLevel
		case "warn":
			level = WarnLevel
		case "error":
			level = ErrorLevel
		}
	}

	return Config{
		Level:     level,
		Output:    os.Stdout,
		AddSource: true,
		JSON:      true,
	}
}

// New creates a new Logger with the given configuration.
func New(cfg Config) *Logger {
	var level slog.Level

	switch cfg.Level {
	case DebugLevel:
		level = slog.LevelDebug
	case InfoLevel:
		level = slog.LevelInfo
	case WarnLevel:
		level = slog.LevelWarn
	case ErrorLevel:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	if cfg.JSON {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// With returns a new Logger with the given attributes added to the context.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// WithContext returns a new context with the logger attached.
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext returns the logger from the context.
// If no logger is found, it returns a default logger.
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return logger
	}

	return &Logger{
		Logger: slog.Default(),
	}
}

// loggerKey is the key used to store the logger in the context.
type loggerKey struct{}

// traceIDKey is the key used to store the trace ID in the context.
type traceIDKey struct{}

// WithTraceID returns a new context with the trace ID attached.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// TraceIDFromContext returns the trace ID from the context.
// If no trace ID is found, it returns an empty string.
func TraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey{}).(string); ok {
		return traceID
	}

	return ""
}

// addTraceIDFromContext adds trace_id field to the log arguments if present in context.
func (l *Logger) addTraceIDFromContext(ctx context.Context, args []any) []any {
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		return append(args, slog.String("trace_id", traceID))
	}

	return args
}

// WithTrace returns a new Logger with trace_id from context added to all logs.
func (l *Logger) WithTrace(ctx context.Context) *Logger {
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		return &Logger{
			Logger: l.Logger.With(slog.String("trace_id", traceID)),
		}
	}

	return l
}

// Context-aware logging methods that automatically add trace_id

// DebugContext logs a debug message with trace_id from context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Logger.Debug(msg, l.addTraceIDFromContext(ctx, args)...)
}

// InfoContext logs an info message with trace_id from context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Logger.Info(msg, l.addTraceIDFromContext(ctx, args)...)
}

// WarnContext logs a warning message with trace_id from context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Logger.Warn(msg, l.addTraceIDFromContext(ctx, args)...)
}

// ErrorContext logs an error message with trace_id from context.
func (l *Logger) ErrorContext(ctx context.Context, err error, msg string, args ...any) {
	l.Logger.Error(msg, l.addTraceIDFromContext(ctx, args)...)
}

// LogRequest logs information about an HTTP request.
func (l *Logger) LogRequest(method, path string, statusCode int, latency time.Duration) {
	l.Info("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", statusCode),
		slog.Duration("latency", latency),
	)
}

// LogRequestContext logs information about an HTTP request with trace_id from context.
func (l *Logger) LogRequestContext(ctx context.Context, method, path string, statusCode int, latency time.Duration) {
	args := []any{
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status", statusCode),
		slog.Duration("latency", latency),
	}
	l.InfoContext(ctx, "request", args...)
}

// LogError logs an error with additional context.
func (l *Logger) LogError(err error, msg string, args ...any) {
	if err != nil {
		l.Error(msg, append(args, slog.Any("error", err))...)
	}
}

// LogErrorContext logs an error with additional context and trace_id from context.
func (l *Logger) LogErrorContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		args = append(args, slog.Any("error", err))
		l.ErrorContext(ctx, err, msg, args...)
	}
}
