package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/repository"
)

type Docs interface {
	GetDocsList(ctx context.Context, s entity.LimitedDocsListInput) ([]entity.Document, error)
	GetDoc(ctx context.Context, docID uuid.UUID, login, token string) (*entity.Document, error)
}

type Service struct {
	Docs
}

func NewService(r *repository.Repository) *Service {
	return &Service{Docs: NewDocsService(r.Docs)}
}
