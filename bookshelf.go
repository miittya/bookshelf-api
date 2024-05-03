package bookshelf

import "errors"

type List struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" validate:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	ID     int
	UserID int
	ListID int
}

type Book struct {
	ID              int    `json:"id" db:"id"`
	Title           string `json:"title" db:"title"`
	Author          string `json:"author" db:"author"`
	Publisher       string `json:"publisher" db:"publisher"`
	PublicationYear int    `json:"publication_year" db:"publication_year"`
	PageCount       int    `json:"page_count" db:"page_count"`
}

type ListsBook struct {
	ID     int
	ListID int
	BookID int
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateBookInput struct {
	Title           *string `json:"title"`
	Author          *string `json:"author"`
	Publisher       *string `json:"publisher"`
	PublicationYear *int    `json:"publication_year"`
	PageCount       *int    `json:"page_count"`
}
