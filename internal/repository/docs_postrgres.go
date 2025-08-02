package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/sirupsen/logrus"
)

type DocsPostgres struct {
	db *sqlx.DB
}

func NewDocsPostgres(db *sqlx.DB) *DocsPostgres {
	return &DocsPostgres{db: db}
}

func (r *DocsPostgres) GetDocsList(ctx *gin.Context, s entity.LimitedDocsListInput) ([]entity.Document, error) {

	queryString := `SELECT 
		d.ID,
		d.FILENAME,
		d.MIME AS MIME,
		d.HAS_FILE AS FILE,
		d.IS_PUBLIC AS PUBLIC,
		d.CREATED_AT AS CREATED    
	FROM DOCUMENTS d `
	//--grant
	args := make([]interface{}, 0)
	argCount := 1

	if s.Key != "" && s.Value != "" {
		queryString += fmt.Sprintf(" WHERE d.%s LIKE $%d ", s.Key, argCount)
		args = append(args, "%"+s.Value+"%")
		argCount++
	}

	queryString += " ORDER BY d.FILENAME, d.CREATED_AT "
	queryString += fmt.Sprintf(" LIMIT $%d ", argCount)
	args = append(args, s.Limit)

	logrus.Debug("queryString=", queryString)
	logrus.Debug("args=", args)

	rows, err := r.db.QueryContext(ctx, queryString, args...)
	if err != nil {
		logrus.Error("DBError:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var docsList []entity.Document
	for rows.Next() {
		var d entity.Document
		if err := rows.Scan(&d.ID, &d.Name, &d.Mime, &d.File, &d.Public, &d.Created); err != nil {
			logrus.Println("Error scanning row:", err)
			continue
		}
		docsList = append(docsList, d)
	}

	logrus.Debug("docs count=", len(docsList))
	return docsList, nil
}

func (r *DocsPostgres) GetDoc(ctx *gin.Context, docID uuid.UUID) (*entity.Document, error) {

	queryString := `
	SELECT id,user_id,filename,path,mime,has_file,is_public,created_at,json_data
                   FROM documents WHERE id=$1 `

	var doc entity.Document
	err := r.db.GetContext(ctx, &doc, queryString, docID)
	if err != nil {
		return nil, err
	}
	rows, _ := r.db.QueryContext(ctx,
		`SELECT login FROM document_grants WHERE doc_id=$1`, docID)
	for rows.Next() {
		var l string
		_ = rows.Scan(&l)
		doc.Grant = append(doc.Grant, l)
	}

	return &doc, nil
}

func (r *DocsPostgres) CreateDocument(ctx *gin.Context, doc *entity.Document) error {
	tx, err := r.db.Begin()

	queryString := `
	INSERT INTO documents (
		id, 
		user_id, 
		filename,
		path, 
		mime, 
		has_file, 
		is_public, 		
		json_data
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err = tx.ExecContext(ctx, queryString,
		doc.ID,
		//doc.UserID,
		"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		doc.Name,
		doc.Path,
		doc.Mime,
		doc.File,
		doc.Public,
		//	pq.Array(doc.Grant),
		doc.JSONData,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	queryString = `
		INSERT INTO document_grants (doc_id, login)
		VALUES ($1, (SELECT login FROM users WHERE id = $2))
		ON CONFLICT (doc_id, login) DO NOTHING`

	_, err = tx.ExecContext(ctx, queryString,
		doc.ID,
		//doc.UserID,
		"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}
