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

func (r *AuthPostgres) CreateUser(login, password string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO users (id, login, password) VALUES ($1,$2,$3) RETURNING id")
	row := r.db.QueryRow(query, uuid.New(), login, password)
	if err := row.Scan(&id); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, login, password, token FROM users where login=$1 AND password=$2")
	err := r.db.Get(&user, query, login, password)

	return &user, err
}

func (r *AuthPostgres) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, login, password, token FROM users where login=$1")
	err := r.db.Get(&user, query, login)

	return &user, err
}

func (r *AuthPostgres) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, login, password, token FROM users where id=$1")
	err := r.db.Get(&user, query, id)

	return &user, err
}

func (r *AuthPostgres) UpdateUserToken(uuid uuid.UUID, token string) error {
	query := fmt.Sprintf("UPDATE users set token=$1 where id=$2")
	_, err := r.db.Exec(query, token, uuid)
	return err
}
