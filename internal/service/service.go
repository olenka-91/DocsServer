package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/olenka-91/DocsServer/internal/repository"
	"github.com/olenka-91/DocsServer/internal/storage"
)

type Docs interface {
	GetDocsList(ctx *gin.Context, s entity.LimitedDocsListInput) ([]entity.Document, error)
	GetDoc(ctx *gin.Context, docID uuid.UUID, token string) (*entity.Document, error)
	PostDoc(ctx *gin.Context, login, token string, meta entity.UploadMeta,
		jsonData entity.JSONB, fileHeader *multipart.FileHeader) (*entity.Document, error)
}

type Service struct {
	Docs
}

func NewService(r *repository.Repository, fs *storage.FileStorage) *Service {
	return &Service{Docs: NewDocsService(r.Docs, fs)}
}
