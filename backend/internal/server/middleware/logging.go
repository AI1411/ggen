package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"g_gen/internal/infra/logger"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewLogging(appLogger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		endpoint := c.Request.RequestURI

		// Skip logging for metrics endpoint
		if endpoint == "/metrics" {
			c.Next()
			return
		}

		// Get or generate trace ID
		var traceID string
		traceID = c.GetHeader("X-Trace-ID")

		if traceID == "" {
			traceID = getTraceID(c.Request.Context())
		}

		// Set trace ID in context
		ctx := logger.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// Set trace ID in Gin context for easy access
		c.Set("trace_id", traceID)

		// Read and log request body
		bufBody, _ := io.ReadAll(c.Request.Body)
		reqBody := map[string]any{}
		_ = json.Unmarshal(bufBody, &reqBody)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bufBody))

		// Log request with trace ID from context
		appLogger.InfoContext(ctx, "request",
			"http_method", c.Request.Method,
			"endpoint", endpoint,
			"header", c.Request.Header,
			"body", reqBody,
		)

		// Wrap response writer to capture response body
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log response with trace ID from context
		appLogger.InfoContext(ctx, "response",
			"endpoint", endpoint,
			"header", writer.Header(),
			"http_status", writer.Status(),
			"body", writer.body.String(),
			"latency_ms", latency.Milliseconds(),
		)
	}
}

func getTraceID(ctx context.Context) string {
	// First check if trace ID already exists in context
	if traceID := logger.TraceIDFromContext(ctx); traceID != "" {
		return traceID
	}

	// Generate new trace ID
	tid, _ := uuid.NewRandom()

	return tid.String()
}
