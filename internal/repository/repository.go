package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
)

type Docs interface {
	GetDocsList(ctx context.Context, s entity.LimitedDocsListInput) ([]entity.Document, error)
	GetDoc(ctx context.Context, docID uuid.UUID) (*entity.Document, error)
}

type Repository struct {
	Docs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Docs: NewDocsPostgres(db)}
}
