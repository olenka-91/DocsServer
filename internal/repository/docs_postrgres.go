package repository

import (
	"fmt"

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

func (r *DocsPostgres) GetDocsList(s entity.LimitedDocsListInput) ([]entity.Document, error) {

	queryString := `SELECT 
		d.ID,
		d.FILENAME,
		d.MIME AS MIME,
		d.HAS_FILE AS FILE,
		d.IS_PUBLIC AS PUBLIC,
		d.CREATED_AT AS CREATED    
	FROM DOCUMENTS d `

	args := make([]interface{}, 0)
	argCount := 1

	if s.Key != "" && s.Value != "" {
		queryString += fmt.Sprintf(" WHERE d.%s LIKE $%d ", groupTable, argCount)
		args = append(args, "%"+s.Value+"%")
		argCount++
	}

	queryString += " ORDER BY d.FILENAME, d.CREATED_AT "
	queryString += fmt.Sprintf(" LIMIT $%d ", argCount)
	args = append(args, s.Limit)

	logrus.Debug("queryString=", queryString)
	logrus.Debug("args=", args)

	rows, err := r.db.Query(queryString, args...)
	if err != nil {
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
