package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	"github.com/zale144/ports/portdomainservice/internal/config"
)

type DB struct {
	*sqlx.DB
}

func Init(cfg *config.Config) (*DB, error) {

	db, err := sqlx.Open("postgres", cfg.DBConnString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db }, nil
}
