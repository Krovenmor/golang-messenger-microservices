package service

import "errors"

var (
	ErrUnknown  = errors.New("unknown error")
	ErrInteranl = errors.New("internal")

	ErrNotEnoughSpace = errors.New("not enough space")
	ErrNotFound       = errors.New("not found")

	ErrAlreadyExists = errors.New("already exists")
)
