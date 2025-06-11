package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"g_gen/internal/infra/logger"
)

const traceIDKey = "trace_id"

func GetTraceID(c *gin.Context) string {
	traceID, exists := c.Get(traceIDKey)
	if !exists {
		return uuid.NewString()
	}

	return traceID.(string)
}

type EmptyResponse struct{}

func handleError(c *gin.Context, err error, appLogger *logger.Logger, message string) {
	res := CreateErrResponse(err)
	res.outputErrorLog(appLogger, message, GetTraceID(c))
	c.AbortWithStatusJSON(res.status, res)
}

func handleValidationError(c *gin.Context, err error, appLogger *logger.Logger, message string) {
	res := createValidateErrorResponse(err)
	res.outputErrorLog(appLogger, message, GetTraceID(c))
	c.AbortWithStatusJSON(res.status, res)
}
