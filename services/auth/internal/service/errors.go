package service

import "errors"

var (
	ErrBadData  = errors.New("bad data")
	ErrInternal = errors.New("internal")
)
