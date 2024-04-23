package postgres

import (
	bookshelf "bookshelf-api"
	"database/sql"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (s *AuthPostgres) CreateUser(user bookshelf.User) (int, error) {
	var id int
	query := "INSERT INTO users(username, password_hash) values ($1, $2) RETURNING id"
	row := s.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthPostgres) GetUser(username, password string) (bookshelf.User, error) {
	var user bookshelf.User
	query := "SELECT id FROM users WHERE username=$1 AND password_hash=$2"
	err := s.db.QueryRow(query, username, password).Scan(&user.ID)
	return user, err
}
