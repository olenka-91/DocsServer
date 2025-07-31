package service

import "errors"

var (
	ErrBadRequest           = errors.New("bad request")           //http.StatusBadRequest = 400
	ErrUnauthorized         = errors.New("unauthorized")          //http.StatusUnauthorized = 401
	ErrForbidden            = errors.New("forbidden")             //http.StatusForbidden = 403
	ErrNotFound             = errors.New("doc not found")         //http.StatusNotFound = 405
	ErrInternalServerError  = errors.New("internal server error") //http.StatusInternalServerError = 500
	ErrMethodNotImplemented = errors.New("not implemented")       //http.StatusMethodNotImplemented = 501
)
