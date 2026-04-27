package apperror

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// AppError is a structured error with an HTTP status and a client-safe message.
// Anything that isn't an AppError is treated as a 500 with a generic message —
// the underlying error is logged but never sent to the client.
type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func Wrap(err error, status int, message string) *AppError {
	return &AppError{Status: status, Message: message, Err: err}
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

// FromGorm maps GORM sentinel errors to AppError. Returns the original error
// untouched if it isn't a known sentinel (caller decides how to handle it).
func FromGorm(err error, notFoundMessage string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(notFoundMessage)
	}
	return err
}

// As is a convenience for unwrapping AppError from the controller.
func As(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
