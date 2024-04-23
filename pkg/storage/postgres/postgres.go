package postgres

import (
	"bookshelf-api/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func New(cfg config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, err
}
