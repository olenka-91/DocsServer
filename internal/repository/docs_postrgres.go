package repository

import (
	"database/sql"
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

func (r *DocsPostgres) GetGrantListByDocID(ctx *gin.Context, docID uuid.UUID) ([]string, error) {
	queryString := `SELECT u.login FROM users u
				INNER JOIN document_grants g ON g.user_id = u.id
				WHERE g.doc_id=$1	
	`
	rows, err := r.db.QueryContext(ctx, queryString, docID)
	if err != nil {
		logrus.Error("DBError:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var grantList []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			logrus.Println("Error scanning row:", err)
			continue
		}
		grantList = append(grantList, s)
	}

	return grantList, nil
}

func (r *DocsPostgres) GetLoginByUserID(ctx *gin.Context, userID uuid.UUID) string {
	queryString := `SELECT login FROM users u				
				WHERE u.id = $1					
	`
	row := r.db.QueryRowContext(ctx, queryString, userID)

	var s string
	if err := row.Scan(&s); err != nil {
		logrus.Println("Error scanning row:", err)
		return ""
	}

	return s
}

func (r *DocsPostgres) GetUserIDByLogin(ctx *gin.Context, login string) uuid.UUID {
	queryString := `SELECT u.id FROM users u				
				WHERE u.login = $1					
	`
	logrus.Debug("login=", login)
	row := r.db.QueryRowContext(ctx, queryString, login)

	var s uuid.UUID
	if err := row.Scan(&s); err != nil {
		logrus.Println("Error scanning row:", err)
		return uuid.Nil
	}

	return s
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
		d.Grant, err = r.GetGrantListByDocID(ctx, d.ID)
		if err != nil {
			logrus.Println("Error getting grant:", err)
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

	doc.Grant, err = r.GetGrantListByDocID(ctx, doc.ID)
	if err != nil {
		return nil, err
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.ExecContext(ctx, queryString,
		doc.ID,
		doc.UserID,
		doc.Name,
		doc.Path,
		doc.Mime,
		doc.File,
		doc.Public,
		doc.JSONData,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	queryString = `
		INSERT INTO document_grants (doc_id, user_id)
		VALUES ($1,  $2)
		ON CONFLICT (doc_id, user_id) DO NOTHING`

	_, err = tx.ExecContext(ctx, queryString,
		doc.ID,
		doc.UserID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	queryString = `
		INSERT INTO document_grants (doc_id, user_id)
		SELECT $1, u.id
		FROM users u 
		WHERE u.login = $2
		ON CONFLICT (doc_id, user_id) DO NOTHING`

	// Подготавливаем запрос один раз
	stmt, err := tx.PrepareContext(ctx, queryString)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare grant statement: %w", err)
	}
	defer stmt.Close()

	// Выполняем для каждого логина
	for _, grantLogin := range doc.Grant {
		_, err = stmt.ExecContext(ctx, doc.ID, grantLogin)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to grant access for %s: %w", grantLogin, err)
		}
		logrus.Debugf("Granted access to %s", grantLogin)
	}
	tx.Commit()
	return err
}

func (r *DocsPostgres) DeleteDoc(ctx *gin.Context, docID uuid.UUID) (bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		"DELETE FROM documents WHERE id = $1",
		docID,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, sql.ErrNoRows
	}

	tx.Commit()
	return true, nil
}
