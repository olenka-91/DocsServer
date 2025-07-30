package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olenka-91/DocsServer/internal/config"
	log "github.com/sirupsen/logrus"
)

const (
	songTable  = "songs"
	groupTable = "groups"
)

func NewPostgresDB(cfg config.Config) (*sqlx.DB, error) {
	log.Debugf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, err
}
