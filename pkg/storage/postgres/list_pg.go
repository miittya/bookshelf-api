package postgres

import (
	bookshelf "bookshelf-api"
	"database/sql"
)

type ListPostgres struct {
	db *sql.DB
}

func NewListPostgres(db *sql.DB) *ListPostgres {
	return &ListPostgres{db: db}
}

func (s *ListPostgres) Create(userID int, list bookshelf.List) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	listsQuery := "INSERT INTO lists(title, description) VALUES ($1, $2) RETURNING id"
	row := tx.QueryRow(listsQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	usersListsQuery := "INSERT INTO users_lists(user_id, list_id) VALUES ($1, $2)"
	_, err = tx.Exec(usersListsQuery, userID, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (s *ListPostgres) GetAll(userID int) ([]bookshelf.List, error) {
	var lists []bookshelf.List
	query := "SELECT l.id, l.title, l.description FROM lists l INNER JOIN users_lists ul ON l.id=ul.list_id WHERE ul.user_id=$1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list bookshelf.List
		err := rows.Scan(&list.ID, &list.Title, &list.Description)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, err
}

func (s *ListPostgres) GetByID(userID, listID int) (bookshelf.List, error) {
	var list bookshelf.List
	query := "SELECT l.id, l.title, l.description FROM lists l INNER JOIN users_lists ul ON l.id=ul.list_id WHERE ul.user_id=$1 AND ul.list_id=$2"
	row := s.db.QueryRow(query, userID, listID)
	err := row.Scan(&list.ID, &list.Title, &list.Description)
	if err != nil {
		return bookshelf.List{}, err
	}
	return list, nil
}

func (s *ListPostgres) Update(userID, listID int, list bookshelf.List, input bookshelf.UpdateListInput) error {
	if input.Title == nil {
		input.Title = &list.Title
	}
	if input.Description == nil {
		input.Description = &list.Description
	}
	query := "UPDATE lists l SET title = $1, description = $2 FROM users_lists ul WHERE l.id = ul.list_id AND ul.list_id = $3 AND ul.user_id = $4"
	_, err := s.db.Exec(query, input.Title, input.Description, listID, userID)
	return err
}

func (s *ListPostgres) Delete(userID, listID int) error {
	query := "DELETE FROM lists l USING users_lists ul WHERE l.id=ul.list_id AND ul.user_id=$1 AND ul.list_id=$2"
	_, err := s.db.Exec(query, userID, listID)
	return err
}
