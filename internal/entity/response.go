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

type RegisterUserResponse struct {
	Login string `json:"login"`
}

type AuthUserResponse struct {
	Token string `json:"token"`
}
