package service

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/repository"
	"github.com/olenka-91/DocsServer/internal/storage"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	defLimit = 10
)

type DocsService struct {
	repo    repository.Docs
	storage *storage.FileStorage
}

func NewDocsService(r repository.Docs, fs *storage.FileStorage) *DocsService {
	return &DocsService{repo: r, storage: fs}
}

func (s *DocsService) GetDocsList(ctx *gin.Context, input entity.LimitedDocsListInput) ([]entity.Document, error) {
	log.Debugf("Fetching list of docs with limit: %+v", input)

	limit := 10
	if input.Limit <= 0 {
		input.Limit = limit
	}

	docs, err := s.repo.GetDocsList(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, ErrNotFound
	}

	if input.Login == "" {
		return nil, ErrUnauthorized
	}

	return docs, nil
}

func (s *DocsService) GetDoc(ctx *gin.Context, docID uuid.UUID, login string) (*entity.Document, error) {
	log.Debugf("Fetching doc with ID: %+v", docID)
	doc, err := s.repo.GetDoc(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, ErrNotFound
	}

	if login == "" {
		return nil, ErrUnauthorized
	}

	if !s.canAccess(ctx, doc, login) {
		return nil, ErrForbidden
	}

	if doc.File && (ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodHead) {

		if err := s.storage.ServeFile(ctx, doc); err != nil {
			return nil, err
		}
		ctx.Abort()
		return nil, nil
	}

	return doc, nil

}

func (s *DocsService) canAccess(ctx *gin.Context, doc *entity.Document, login string) bool {
	if doc.Public {
		return true
	}
	if login == s.repo.GetLoginByUserID(ctx, doc.UserID) {
		return true
	}
	for _, l := range doc.Grant {
		if l == login {
			return true
		}
	}
	return false
}

func (s *DocsService) PostDoc(ctx *gin.Context, login string, meta entity.UploadMeta,
	jsonData entity.JSONB, fileHeader *multipart.FileHeader) (*entity.Document, error) {
	logrus.Debugf("Posting doc to storage.")

	if meta.File && fileHeader == nil {
		return nil, ErrBadRequest
	}

	userID := s.repo.GetUserIDByLogin(ctx, login)
	doc := entity.Document{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     meta.Name,
		Mime:     meta.Mime,
		File:     meta.File,
		Public:   meta.Public,
		Grant:    meta.Grant,
		JSONData: jsonData,
	}

	var storedMime, filePath string
	if doc.File {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, storedMime, filePath, err = s.storage.SaveFile(doc.ID, file, meta.Name)
		if err != nil {
			logrus.Errorf("Failed to save file: %v", err)
			return nil, err
		}

		logrus.Debugf("Doc saved at:%s", filePath)
	}

	if storedMime != "" {
		doc.Mime = storedMime
	}

	if filePath != "" {
		doc.Path = filePath
	}

	if err := s.repo.CreateDocument(ctx, &doc); err != nil {
		logrus.Errorf("Failed to create document: %v", err)
		s.storage.DeleteFile(doc.ID, doc.Name)
		return nil, err
	}

	return &doc, nil
}

func (s *DocsService) DeleteDoc(ctx *gin.Context, docID uuid.UUID, login string) (*entity.DelResponse, error) {
	log.Debugf("Deleting doc with ID: %+v", docID)

	if login == "" {
		return nil, ErrUnauthorized
	}

	doc, err := s.repo.GetDoc(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, ErrNotFound
	}

	if login != s.repo.GetLoginByUserID(ctx, doc.UserID) {
		return nil, ErrForbidden
	}

	if doc.File {
		if err := s.storage.DeleteFile(docID, doc.Name); err != nil {
			return nil, err
		}
	}

	if _, err := s.repo.DeleteDoc(ctx, docID); err != nil {
		return nil, err
	}

	return &entity.DelResponse{docID: true}, nil

}
