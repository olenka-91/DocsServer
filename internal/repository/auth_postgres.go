package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(login, password string) (string, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO users (id, login, password) VALUES ($1,$2,$3) RETURNING id")
	row := r.db.QueryRow(query, uuid.New(), login, password)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return login, nil
}

func (r *AuthPostgres) GetUser(login, password string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, login, password FROM users where login=$1 AND password=$2")
	err := r.db.Get(&user, query, login, password)

	return user, err
}
