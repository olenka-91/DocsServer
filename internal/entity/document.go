package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type LimitedDocsListInput struct {
	Token string `json:"token"`
	Login string `json:"login"` //опционально — если не указан — то список своих
	Key   string `json:"key"`   //имя колонки для фильтрации
	Value string `json:"value"` //- значение фильтра
	Limit int    `json:"limit"` //кол-во документов в списке
}

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

type Document struct {
	ID       uuid.UUID    `db:"id"          json:"id"`
	OwnerID  uuid.UUID    `db:"owner_id"    json:"-"`
	Name     string       `db:"filename"    json:"filename"`
	Path     string       `db:"path"        json:"-"`
	Mime     string       `db:"mime"        json:"mime"`
	File     bool         `db:"has_file"    json:"file"`
	Public   bool         `db:"is_public"   json:"public"`
	Created  sql.NullTime `db:"created_at"  json:"created"`
	Grant    []string     `db:"grant"       json:"grant,omitempty"`
	JSONData JSONB        `db:"json_data"   json:"json,omitempty"`
}

// CREATE TABLE DOCUMENTS (
//     ID         UUID PRIMARY KEY,
//     OWNER_ID   UUID REFERENCES USERS(ID),
//     FILENAME   TEXT NOT NULL,
//     PATH       TEXT NOT NULL,
//     CREATED_AT TIMESTAMPTZ DEFAULT NOW()
//     );
// 	ALTER TABLE DOCUMENTS
// 	ADD COLUMN MIME  TEXT NOT NULL,
// 	ADD COLUMN HAS_FILE    BOOLEAN NOT NULL,
// 	ADD COLUMN IS_PUBLIC   BOOLEAN NOT NULL;

type DocsData struct {
	Docs []Document `json:"docs"`
}

type DocsResponse struct {
	Data DocsData `json:"data"`
}
