package errors

import (
	"errors"
	"fmt"
)

const (
	CodeSuccess         = 0
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeInternalError   = 500
	CodeUserNotFound    = 1001
	CodeInvalidPassword = 1002
	CodeInsufficientStock = 1004
	CodeDuplicateEntry  = 1005
	CodeRecordNotFound  = 1006
)

type AppError struct {
	Code    int
	Message string
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func Wrap(err error, code int, message string) error {
	return fmt.Errorf("%w", NewAppError(code, message))
}

var (
	ErrBadRequest    = NewAppError(CodeBadRequest, "bad request")
	ErrUnauthorized  = NewAppError(CodeUnauthorized, "unauthorized")
	ErrForbidden     = NewAppError(CodeForbidden, "forbidden")
	ErrNotFound      = NewAppError(CodeNotFound, "not found")
	ErrInternalError = NewAppError(CodeInternalError, "internal server error")
)
