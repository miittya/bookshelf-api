package bookshelf

type User struct {
	ID       int    `json:"-" db:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
