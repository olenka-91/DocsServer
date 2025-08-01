package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
)

type Docs interface {
	GetDocsList(ctx *gin.Context, s entity.LimitedDocsListInput) ([]entity.Document, error)
	GetDoc(ctx *gin.Context, docID uuid.UUID) (*entity.Document, error)
	CreateDocument(ctx *gin.Context, doc *entity.Document) error
}

type Repository struct {
	Docs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Docs: NewDocsPostgres(db)}
}
