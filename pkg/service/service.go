package service

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/storage"
)

type Authorization interface {
	CreateUser(user bookshelf.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type List interface {
	Create(userID int, list bookshelf.List) (int, error)
	GetAll(userID int) ([]bookshelf.List, error)
	GetByID(userID, listID int) (bookshelf.List, error)
	Update(userID, listID int, input bookshelf.UpdateListInput) error
	Delete(userID, listID int) error
}

type Book interface {
	Create(userID, listID int, book bookshelf.Book) (int, error)
	GetAll(userID, listID int) ([]bookshelf.Book, error)
	GetByID(userID, bookID int) (bookshelf.Book, error)
	Update(userID, bookID int, input bookshelf.UpdateBookInput) error
	Delete(userID, bookID int) error
}
type Service struct {
	Authorization
	List
	Book
}

func New(storage *storage.Storage) *Service {
	return &Service{
		Authorization: NewAuthService(storage.Authorization),
		List:          NewListService(storage.List),
		Book:          NewBookService(storage.Book, storage.List),
	}
}
