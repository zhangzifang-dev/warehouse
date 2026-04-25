package errors

import (
	"errors"
	"testing"
)

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"CodeSuccess", CodeSuccess, 0},
		{"CodeBadRequest", CodeBadRequest, 400},
		{"CodeUnauthorized", CodeUnauthorized, 401},
		{"CodeForbidden", CodeForbidden, 403},
		{"CodeNotFound", CodeNotFound, 404},
		{"CodeInternalError", CodeInternalError, 500},
		{"CodeUserNotFound", CodeUserNotFound, 1001},
		{"CodeInvalidPassword", CodeInvalidPassword, 1002},
		{"CodeInsufficientStock", CodeInsufficientStock, 1004},
		{"CodeDuplicateEntry", CodeDuplicateEntry, 1005},
		{"CodeRecordNotFound", CodeRecordNotFound, 1006},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.expected)
			}
		})
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError(CodeNotFound, "resource not found")

	if err.Code != CodeNotFound {
		t.Errorf("Code = %d, want %d", err.Code, CodeNotFound)
	}
	if err.Message != "resource not found" {
		t.Errorf("Message = %s, want %s", err.Message, "resource not found")
	}
}

func TestAppErrorError(t *testing.T) {
	err := NewAppError(CodeBadRequest, "invalid input")

	if err.Error() != "invalid input" {
		t.Errorf("Error() = %s, want %s", err.Error(), "invalid input")
	}
}

func TestIsAppError(t *testing.T) {
	appErr := NewAppError(CodeNotFound, "not found")
	regularErr := errors.New("regular error")

	if !IsAppError(appErr) {
		t.Error("IsAppError(appErr) = false, want true")
	}
	if IsAppError(regularErr) {
		t.Error("IsAppError(regularErr) = true, want false")
	}
}

func TestGetAppError(t *testing.T) {
	appErr := NewAppError(CodeNotFound, "not found")
	regularErr := errors.New("regular error")

	result, ok := GetAppError(appErr)
	if !ok {
		t.Error("GetAppError(appErr) ok = false, want true")
	}
	if result.Code != CodeNotFound {
		t.Errorf("GetAppError(appErr).Code = %d, want %d", result.Code, CodeNotFound)
	}

	result, ok = GetAppError(regularErr)
	if ok {
		t.Error("GetAppError(regularErr) ok = true, want false")
	}
	if result != nil {
		t.Errorf("GetAppError(regularErr) = %v, want nil", result)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		code     int
		message  string
	}{
		{"ErrBadRequest", ErrBadRequest, CodeBadRequest, "bad request"},
		{"ErrUnauthorized", ErrUnauthorized, CodeUnauthorized, "unauthorized"},
		{"ErrForbidden", ErrForbidden, CodeForbidden, "forbidden"},
		{"ErrNotFound", ErrNotFound, CodeNotFound, "not found"},
		{"ErrInternalError", ErrInternalError, CodeInternalError, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("%s.Code = %d, want %d", tt.name, tt.err.Code, tt.code)
			}
			if tt.err.Message != tt.message {
				t.Errorf("%s.Message = %s, want %s", tt.name, tt.err.Message, tt.message)
			}
		})
	}
}
