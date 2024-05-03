package storage

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/storage/postgres"
	"database/sql"
)

type Authorization interface {
	CreateUser(user bookshelf.User) (int, error)
	GetUser(username, password string) (bookshelf.User, error)
}

type List interface {
	Create(userID int, list bookshelf.List) (int, error)
	GetAll(userID int) ([]bookshelf.List, error)
	GetByID(userID, listID int) (bookshelf.List, error)
	Update(userID, listID int, list bookshelf.List, input bookshelf.UpdateListInput) error
	Delete(userID, listID int) error
}

type Book interface {
	Create(listID int, book bookshelf.Book) (int, error)
	GetAll(userID, listID int) ([]bookshelf.Book, error)
	GetByID(userID, bookID int) (bookshelf.Book, error)
	Update(userID, bookID int, input bookshelf.UpdateBookInput) error
	Delete(userID, bookID int) error
}

type Storage struct {
	Authorization
	List
	Book
}

func New(db *sql.DB) *Storage {
	return &Storage{
		Authorization: postgres.NewAuthPostgres(db),
		List:          postgres.NewListPostgres(db),
		Book:          postgres.NewBookPostgres(db),
	}
}
