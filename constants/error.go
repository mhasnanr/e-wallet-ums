package constants

import "net/http"

var (
	ErrFieldBadRequest  = "data is not valid"
	ErrDuplicateEmail   = "email is already registered"
	ErrUserNotFound     = "user not found"
	ErrRequiredEmail    = "email is required"
	ErrRequiredUsername = "username is required"
	ErrRequiredPassword = "password is required"
	ErrSessionNotFound 	= "user session not found"
)

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

var (
	ErrorDuplicateEmail  = NewAppError(http.StatusConflict, ErrDuplicateEmail)
	ErrorRequiredEmail   = NewAppError(http.StatusBadRequest, ErrRequiredEmail)
	ErrorUserNotFound    = NewAppError(http.StatusNotFound, ErrUserNotFound)
	ErrorSessionNotFound = NewAppError(http.StatusNotFound, ErrSessionNotFound)
)
