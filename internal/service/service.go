package service

import (
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/repository"
)

type Docs interface {
	GetDocsList(s entity.LimitedDocsListInput) ([]entity.Document, error)
}

type Service struct {
	Docs
}

func NewService(r *repository.Repository) *Service {
	return &Service{Docs: NewDocsService(r.Docs)}
}
