package myerrors

import "github.com/pkg/errors"

type (
	ErrorCode    string
	ErrorMessage string
)

const (
	SystemError             ErrorCode = "E100000" // システムエラー
	ValidationError         ErrorCode = "E100001"
	PrefectureNotFoundError ErrorCode = "E100002" // 都道府県が存在しないエラー
)

const (
	SystemErrorMessage             ErrorMessage = "システムエラーが発生しました"
	ValidationErrorMessage         ErrorMessage = "入力値に誤りがあります"
	PrefectureNotFoundErrorMessage ErrorMessage = "都道府県は存在しません"
)

func NewAPIError(code ErrorCode, msg ErrorMessage, originalErr error, internalMsg string) *APIError {
	return &APIError{
		Code:     code,
		Message:  msg,
		internal: errors.Wrap(originalErr, internalMsg),
	}
}

type APIError struct {
	Code     ErrorCode
	Message  ErrorMessage
	internal error
}

func (e APIError) Error() string {
	if e.internal != nil {
		return e.internal.Error()
	}
	return ""
}

func (e APIError) ErrorCode() ErrorCode {
	return e.Code
}

func (e APIError) ErrorMessage() ErrorMessage {
	return e.Message
}
