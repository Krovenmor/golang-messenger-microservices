package service

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnknown       = errors.New("unknown")
	ErrInternal      = errors.New("internal")
	ErrAlreadyExists = errors.New("already exists")
)
