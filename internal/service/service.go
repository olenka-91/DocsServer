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
	SignUp(name, password string) (map[string]string, error)
	SignIn(name, password string) (map[string]string, error)
	RefreshToken(refreshToken string) (map[string]string, error)
	Logout(userID uuid.UUID) error
}

type Service struct {
	Docs
	Authorization
}

func NewService(r *repository.Repository, fs *storage.FileStorage) *Service {
	return &Service{Docs: NewDocsService(r.Docs, fs), Authorization: NewAuthService(r.Authorization)}
}
