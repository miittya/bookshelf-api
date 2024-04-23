package service

import (
	bookshelf "bookshelf-api"
	"bookshelf-api/pkg/storage"
)

type ListService struct {
	storage storage.List
}

func NewListService(storage storage.List) *ListService {
	return &ListService{storage: storage}
}

func (s *ListService) Create(userID int, list bookshelf.List) (int, error) {
	return s.storage.Create(userID, list)
}

func (s *ListService) GetAll(userID int) ([]bookshelf.List, error) {
	return s.storage.GetAll(userID)
}

func (s *ListService) GetByID(userID, listID int) (bookshelf.List, error) {
	return s.storage.GetByID(userID, listID)
}

func (s *ListService) Update(userID, listID int, input bookshelf.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.storage.Update(userID, listID, input)
}

func (s *ListService) Delete(userID, listID int) error {
	return s.storage.Delete(userID, listID)
}
