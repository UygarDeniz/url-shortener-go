package errors

import "errors"

var (
	ErrDatabaseAccess     = errors.New("database access error")
	ErrInvalidURL         = errors.New("invalid URL")
	ErrDuplicateShortCode = errors.New("short code already exists")
)

var (
	ErrInvalidDataFormat   = errors.New("invalid data format")
	ErrShortCodeGeneration = errors.New("short code generation error")
)
