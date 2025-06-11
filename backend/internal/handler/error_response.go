package handler

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	myerrors "g_gen/internal/errors"
	"g_gen/internal/infra/logger"
)

type ErrorResponse struct {
	// 内部のエラーコード
	Code myerrors.ErrorCode `json:"code" example:"E100000"`
	// 内部のエラーメッセージ
	Message myerrors.ErrorMessage `json:"message" example:"システムエラーが発生しました。"`
	err     error
	// HTTPステータス
	status int
}

func (r *ErrorResponse) outputErrorLog(appLogger *logger.Logger, message, traceID string) {
	msg := fmt.Sprintf("%s: %s", message, r.err.Error())
	if r.status == http.StatusInternalServerError {
		appLogger.Error(msg, zap.Error(r.err), zap.String("trace_id", traceID))
	} else {
		appLogger.Debug(msg, zap.Error(r.err), zap.String("trace_id", traceID))
	}
}

type ValidationError struct {
	Attribute string `json:"attribute" example:"name"`
	Tag       string `json:"tag" example:"required"`
	Message   string `json:"message" example:"項目は必須です"`
}

type ErrorResponseDetail struct {
	Code    myerrors.ErrorCode    `json:"code" example:"E110000"`
	Message myerrors.ErrorMessage `json:"message" example:"パラメータが正しく設定されていません。"`
	Details []ValidationError     `json:"details,omitempty"`
	err     error
	status  int
}

func (r *ErrorResponseDetail) outputErrorLog(appLogger *logger.Logger, message, traceID string) {
	msg := message
	if r.status == http.StatusInternalServerError {
		appLogger.Error(msg, zap.Error(r.err), zap.String("trace_id", traceID))
	} else {
		appLogger.Debug(msg, zap.Error(r.err), zap.String("trace_id", traceID))
	}
}

func CreateErrResponse(err error) *ErrorResponse {
	switch cErr := err.(type) {
	case *myerrors.APIError:
		switch cErr.Code {
		case myerrors.SystemError:
			return &ErrorResponse{
				Code:    cErr.Code,
				Message: cErr.Message,
				err:     cErr,
				status:  http.StatusInternalServerError,
			}
		case myerrors.ValidationError:
			return &ErrorResponse{
				Code:    cErr.Code,
				Message: cErr.Message,
				err:     cErr,
				status:  http.StatusBadRequest,
			}
		case myerrors.PrefectureNotFoundError:
			return &ErrorResponse{
				Code:    cErr.Code,
				Message: cErr.Message,
				err:     cErr,
				status:  http.StatusNotFound,
			}
		default:
			return &ErrorResponse{
				Code:    cErr.Code,
				Message: cErr.Message,
				err:     cErr,
				status:  http.StatusInternalServerError,
			}
		}
	default:
		return &ErrorResponse{
			Code:    myerrors.SystemError,
			Message: myerrors.SystemErrorMessage,
			err:     err,
			status:  http.StatusInternalServerError,
		}
	}
}
