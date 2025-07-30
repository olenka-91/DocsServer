package service

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrAccessDenied = errors.New("access denied")
	ErrInvalidInput = errors.New("invalid input parameters")
)
