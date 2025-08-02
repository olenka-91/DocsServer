package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"-" db:"id"`
	Login    string    `json:"login" db:"login" binding:"required"`
	Password string    `json:"password" db:"password" binding:"required"`
}
