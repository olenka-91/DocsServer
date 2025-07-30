package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
)

type Docs interface {
	GetDocsList(s entity.LimitedDocsListInput) ([]entity.Document, error)
}

type Repository struct {
	Docs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Docs: NewDocsPostgres(db)}
}
