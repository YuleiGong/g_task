package server

import "errors"

var (
	ErrTimeout = errors.New("exec Func time out")
	ErrFunc    = errors.New("func format error")
)
