package service

import "errors"

var (
	ErrUnknown           = errors.New("unknown")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternal          = errors.New("internal")
	ErrConnNormalClosure = errors.New("normal closure")
)
