package service

import (
	"context"

	"github.com/google/uuid"
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

func (r *DocsService) GetDocsList(ctx context.Context, s entity.LimitedDocsListInput) ([]entity.Document, error) {
	log.Debugf("Fetching list of docs with limit: %+v", s)

	limit := 10
	if s.Limit <= 0 {
		s.Limit = limit
	}

	docs, err := r.repo.GetDocsList(ctx, s)
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, ErrNotFound
	}

	if s.Login == "" {
		return nil, ErrForbidden
	}

	if s.Token == "" {
		return nil, ErrUnauthorized
	}

	// if login == r.user.OwnerLogin(ctx, doc.OwnerID) {
	// 	return doc, nil
	// }

	// for _, l := range doc.Grant {
	// 	if l == login {
	// 		return doc, nil
	// 	}
	// }

	return docs, nil

}

func (r *DocsService) GetDoc(ctx context.Context, docID uuid.UUID, login, token string) (*entity.Document, error) {
	log.Debugf("Fetching doc with ID: %+v", docID)
	doc, err := r.repo.GetDoc(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, ErrNotFound
	}

	if doc.Public {
		return doc, nil
	}

	//if token != "" && token == doc.JSONData["token"] {
	if token != "" {
		return doc, nil
	}

	if login == "" {
		return nil, ErrForbidden
	}

	// if login == r.user.OwnerLogin(ctx, doc.OwnerID) {
	// 	return doc, nil
	// }

	// for _, l := range doc.Grant {
	// 	if l == login {
	// 		return doc, nil
	// 	}
	// }

	return nil, ErrForbidden

}
