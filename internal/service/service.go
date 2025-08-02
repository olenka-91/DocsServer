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
	GetDoc(ctx *gin.Context, docID uuid.UUID, login string) (*entity.Document, error)
	PostDoc(ctx *gin.Context, login string, meta entity.UploadMeta,
		jsonData entity.JSONB, fileHeader *multipart.FileHeader) (*entity.Document, error)
	DeleteDoc(ctx *gin.Context, docID uuid.UUID, login string) (*entity.DelResponse, error)
}

type Authorization interface {
	CreateUser(login, password string) (string, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (string, error)
	ValidateAdminToken(adminToken string) bool
}

type Service struct {
	Docs
	Authorization
}

func NewService(r *repository.Repository, fs *storage.FileStorage, adminToken string) *Service {
	return &Service{Docs: NewDocsService(r.Docs, fs), Authorization: NewAuthService(r.Authorization, adminToken)}
}
