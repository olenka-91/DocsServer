package service

import (
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/repository"
	log "github.com/sirupsen/logrus"
)

const (
	defLimit = 10
)

type DocsService struct {
	repo repository.Docs
}

func NewDocsService(r repository.Docs) *DocsService {
	return &DocsService{repo: r}
}

func (r *DocsService) GetDocsList(s entity.LimitedDocsListInput) ([]entity.Document, error) {
	log.Debugf("Fetching list of docs with limit: %+v", s)

	if s.Limit < 1 {
		log.Warn("Invalid limit, defaulting ", defLimit)
		s.Limit = defLimit
	}

	return r.repo.GetDocsList(s)
}
