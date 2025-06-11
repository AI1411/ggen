package myerrors

import "github.com/pkg/errors"

type ErrorCode string
type ErrorMessage string

const (
	SystemError                            ErrorCode = "E100000" // システムエラー
	ValidationError                        ErrorCode = "E100001" // バリデーションエラー
	PrefectureNotFoundError                ErrorCode = "E100002" // 都道府県が存在しないエラー
	EmailVarificationTokenNotFound         ErrorCode = "E100003" // メール認証トークンが存在しないエラー
	EmailVarificationTokenUsedError        ErrorCode = "E100004" // メール認証トークンが使用済みエラー
	EmailVarificationTokenExpiredError     ErrorCode = "E100005" // メール認証トークンが期限切れエラー
	UserNotFoundError                      ErrorCode = "E100006" // ユーザーが存在しないエラー
	EmailVarificationTokenAlreadyUsedError ErrorCode = "E100007" // メール認証トークンが既に使用済みエラー
)

const (
	SystemErrorMessage                            ErrorMessage = "システムエラーが発生しました"
	ValidationErrorMessage                        ErrorMessage = "入力値に誤りがあります"
	PrefectureNotFoundErrorMessage                ErrorMessage = "都道府県は存在しません"
	EmailVarificationTokenNotFoundErrorMessage    ErrorMessage = "メール認証トークンが存在しません"
	EmailVarificationTokenUsedErrorMessage        ErrorMessage = "メール認証トークンは使用済みです"
	EmailVarificationTokenExpiredErrorMessage     ErrorMessage = "メール認証トークンは期限切れです"
	UserNotFoundErrorMessage                      ErrorMessage = "ユーザーが存在しません"
	EmailVarificationTokenAlreadyUsedErrorMessage ErrorMessage = "メール認証トークンは既に使用済みです"
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
