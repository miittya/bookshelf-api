package service

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/storage"
)

type BookService struct {
	storage     storage.Book
	listStorage storage.List
}

func NewBookService(storage storage.Book, listStorage storage.List) *BookService {
	return &BookService{
		storage:     storage,
		listStorage: listStorage,
	}
}

func (s *BookService) Create(userID, listID int, book bookshelf.Book) (int, error) {
	_, err := s.listStorage.GetByID(userID, listID)
	if err != nil {
		return 0, err
	}

	return s.storage.Create(listID, book)
}

func (s *BookService) GetAll(userID, listID int) ([]bookshelf.Book, error) {
	return s.storage.GetAll(userID, listID)
}

func (s *BookService) GetByID(userID, bookID int) (bookshelf.Book, error) {
	return s.storage.GetByID(userID, bookID)
}

func (s *BookService) Update(userID, bookID int, input bookshelf.UpdateBookInput) error {
	return s.storage.Update(userID, bookID, input)
}

func (s *BookService) Delete(userID, bookID int) error {
	return s.storage.Delete(userID, bookID)
}
