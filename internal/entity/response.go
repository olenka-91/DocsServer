package entity

import "github.com/google/uuid"

type UploadMeta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

type DelResponse map[uuid.UUID]bool

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

/*
type RegisterUserResponse struct {
	Login string `json:"login"`
}

type AuthUserResponse struct {
	Token string `json:"token"`
}

type AuthUserLogoutResponse map[string]bool*/
