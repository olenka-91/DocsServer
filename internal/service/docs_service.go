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
		return nil, ErrForbidden
	}

	if input.Token == "" {
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

func (s *DocsService) GetDoc(ctx *gin.Context, docID uuid.UUID, token string) (*entity.Document, error) {
	log.Debugf("Fetching doc with ID: %+v", docID)
	doc, err := s.repo.GetDoc(ctx, docID)
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, ErrNotFound
	}

	if token == "" {
		return nil, ErrForbidden
	}

	//!doc.Public &&

	if doc.File && (ctx.Request.Method == http.MethodGet || ctx.Request.Method == http.MethodHead) {
		// Внутри ServeFile выставляются заголовки и пишется тело ответа.
		if err := s.storage.ServeFile(ctx, doc); err != nil {
			return nil, err
		}
		// Ответ сформирован, дальнейшая цепочка не нужна.
		ctx.Abort()
		return nil, nil
	}

	// if login == r.user.OwnerLogin(ctx, doc.OwnerID) {
	// 	return doc, nil
	// }

	// for _, l := range doc.Grant {
	// 	if l == login {
	// 		return doc, nil
	// 	}
	// }

	return doc, nil

}

func (s *DocsService) PostDoc(ctx *gin.Context, login, token string, meta entity.UploadMeta,
	jsonData entity.JSONB, fileHeader *multipart.FileHeader) (*entity.Document, error) {

	if meta.File && fileHeader == nil {
		return nil, ErrBadRequest
	}

	// Создаем документ
	doc := entity.Document{
		ID: uuid.New(),
		//	OwnerID:  user.ID,
		Name:     meta.Name,
		Mime:     meta.Mime,
		File:     meta.File,
		Public:   meta.Public,
		Grant:    meta.Grant,
		JSONData: jsonData,
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Сохраняем файл в хранилище
	_, storedMime, err := s.storage.SaveFile(doc.ID, file, meta.Name)
	if err != nil {
		logrus.Errorf("Failed to save file: %v", err)
		return nil, err
	}

	// Обновляем информацию о документе
	if storedMime != "" {
		doc.Mime = storedMime
	}

	// Сохраняем документ в БД
	if err := s.repo.CreateDocument(ctx, &doc); err != nil {
		logrus.Errorf("Failed to create document: %v", err)
		return nil, err
	}

	return &doc, nil
}
