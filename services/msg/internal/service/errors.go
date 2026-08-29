package service

import "errors"

var (
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotFound            = errors.New("not found")
	ErrChatNotFound        = errors.New("chat not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmojiNotFound       = errors.New("emoji not found")
	ErrChatNotFoundOrEmpty = errors.New("chat not found or empty")
	ErrNotFoundOrForbidden = errors.New("forbidden or not found")
	ErrForbidden           = errors.New("forbidden")
	ErrInternal            = errors.New("internal")
	ErrTooMuch             = errors.New("too much")

	ErrUnknown = errors.New("unknown")
)
