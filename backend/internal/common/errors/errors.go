package errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid username or password")
	ErrUsernameTaken   = errors.New("Username already taken")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("admin portal access denied")
	ErrInternalServer  = errors.New("internal server error")
)
