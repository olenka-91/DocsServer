package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
)

type Authorization interface {
	CreateUser(login, password string) (string, error)
	GetUser(username, password string) (entity.User, error)
}

type Docs interface {
	GetDocsList(ctx *gin.Context, s entity.LimitedDocsListInput) ([]entity.Document, error)
	GetDoc(ctx *gin.Context, docID uuid.UUID) (*entity.Document, error)
	CreateDocument(ctx *gin.Context, doc *entity.Document) error
	DeleteDoc(ctx *gin.Context, docID uuid.UUID) (bool, error)
	GetLoginByUserID(ctx *gin.Context, userID uuid.UUID) string
	GetUserIDByLogin(ctx *gin.Context, login string) uuid.UUID
}

type Repository struct {
	Docs
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Docs: NewDocsPostgres(db),
		Authorization: NewAuthPostgres(db)}
}
