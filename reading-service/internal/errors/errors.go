package errors

import "errors"

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrEmptyShortCode = errors.New("short code cannot be empty")
	ErrInvalidInput   = errors.New("invalid input")
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrDatabaseAccess      = errors.New("database access error")
	ErrDataFormat          = errors.New("data format error")
)
