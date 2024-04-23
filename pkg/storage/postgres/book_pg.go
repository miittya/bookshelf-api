package postgres

import (
	bookshelf "bookshelf-api"
	"database/sql"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (s *BookPostgres) Create(listID int, book bookshelf.Book) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var bookID int
	createBookQuery := "INSERT INTO books(title, author, publisher, publication_year, page_count) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	row := tx.QueryRow(createBookQuery, book.Title, book.Author, book.Publisher, book.PublicationYear, book.PageCount)
	err = row.Scan(&bookID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListsBooksQuery := "INSERT INTO lists_books(list_id, book_id) VALUES ($1, $2)"
	_, err = tx.Exec(createListsBooksQuery, listID, bookID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return bookID, tx.Commit()
}

func (s *BookPostgres) GetAll(userID, listID int) ([]bookshelf.Book, error) {
	var books []bookshelf.Book
	query := "SELECT b.title, b.author, b.publisher, b.publication_year, b.page_count FROM books b INNER JOIN lists_books lb ON b.id = lb.book_id INNER JOIN users_lists ul ON lb.list_id = ul.list_id WHERE lb.list_id = $1 AND ul.user_id = $2"
	rows, err := s.db.Query(query, listID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book bookshelf.Book
		err := rows.Scan(&book.Title, &book.Author, &book.Publisher, &book.PublicationYear, &book.PageCount)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
func (s *BookPostgres) GetByID(userID, bookID int) (bookshelf.Book, error) {
	var book bookshelf.Book
	query := "SELECT b.title, b.author, b.publisher, b.publication_year, b.page_count FROM books b INNER JOIN lists_books lb ON b.id = lb.book_id INNER JOIN users_lists ul ON lb.list_id = ul.list_id WHERE b.id = $1 AND ul.user_id = $2"
	row := s.db.QueryRow(query, bookID, userID)
	err := row.Scan(&book.Title, &book.Author, &book.Publisher, &book.PublicationYear, &book.PageCount)
	if err != nil {
		return bookshelf.Book{}, err
	}
	return book, nil
}
func (s *BookPostgres) Update(userID, bookID int, input bookshelf.UpdateBookInput) error {
	book, err := s.GetByID(userID, bookID)
	if err != nil {
		return err
	}
	if input.Title == nil {
		input.Title = &book.Title
	}
	if input.Author == nil {
		input.Author = &book.Author
	}
	if input.Publisher == nil {
		input.Publisher = &book.Publisher
	}
	if input.PublicationYear == nil {
		input.PublicationYear = &book.PublicationYear
	}
	if input.PageCount == nil {
		input.PageCount = &book.PageCount
	}
	query := "UPDATE books b SET title = $1, author = $2, publisher = $3, publication_year = $4, page_count = $5 FROM lists_books lb, users_lists ul WHERE b.id = lb.book_id AND lb.list_id = ul.list_id AND ul.user_id = $6 AND b.id = $7"
	_, err = s.db.Exec(query, input.Title, input.Author, input.Publisher, input.PublicationYear, input.PageCount, userID, bookID)
	return err
}
func (s *BookPostgres) Delete(userID, bookID int) error {
	query := "DELETE FROM books b USING lists_books lb, users_lists ul WHERE  b.id = lb.book_id AND lb.list_id = ul.list_id AND ul.user_id = $1 AND b.id = $2"
	_, err := s.db.Exec(query, userID, bookID)
	return err
}
