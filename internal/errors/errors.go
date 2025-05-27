package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrInternal     = errors.New("internal server error")
	ErrInvalidInput = errors.New("invalid input format")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("resource conflict")
	ErrValidation   = errors.New("validation failed")
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already taken")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
)

var (
	ErrPostNotFound        = errors.New("post not found")
	ErrInvalidPostID       = errors.New("invalid post ID")
	ErrPostTitleRequired   = errors.New("post title is required")
	ErrPostContentRequired = errors.New("post content is required")
	ErrVersionConflict     = errors.New("version conflict - post was modified by another request")
)

var (
	ErrUserIDNotFound = errors.New("user_id not found in context")
)

type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

func (e AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func NewAppError(err error, statusCode int, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrNotFound, http.StatusNotFound, message)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrBadRequest, http.StatusBadRequest, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(ErrInternal, http.StatusInternalServerError, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrUnauthorized, http.StatusUnauthorized, message)
}

func NewValidationError(message string) *AppError {
	return NewAppError(ErrValidation, http.StatusBadRequest, message)
}

func NewConflictError(message string) *AppError {
	return NewAppError(ErrConflict, http.StatusConflict, message)
}
