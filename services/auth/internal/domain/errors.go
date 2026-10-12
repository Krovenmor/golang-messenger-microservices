package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

var (
	ErrInvalidCode        = errors.New("invalid verification code")
	ErrCodeAlreadySent    = errors.New("verification code already sent")
	ErrUserAlreadyExists  = errors.New("login or email already taken")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInvalidToken       = errors.New("invalid refresh token")
	ErrUserNotFound       = errors.New("user not found")
	ErrBanned             = errors.New("account is banned")
	ErrInternal           = errors.New("internal error")
)
