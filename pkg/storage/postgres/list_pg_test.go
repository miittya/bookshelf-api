package postgres

import (
	bookshelf "bookshelf-api"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	list := NewListPostgres(db)

	type args struct {
		userID int
		list   bookshelf.List
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO lists").WithArgs("title", "description").WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: args{
				userID: 1,
				list: bookshelf.List{
					Title:       "title",
					Description: "description",
				},
			},
			want: 1,
		},
		{
			name: "Empty fields",
			mock: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO lists").WithArgs("", "description").WillReturnRows(rows)
				mock.ExpectRollback()
			},
			input: args{
				userID: 1,
				list: bookshelf.List{
					Title:       "",
					Description: "description",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := list.Create(tt.input.userID, tt.input.list)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	list := NewListPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		userID  int
		want    []bookshelf.List
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM lists l INNER JOIN users_lists ul ON (.+) WHERE (.+)").WithArgs(1).WillReturnRows(rows)
			},
			userID: 1,
			want: []bookshelf.List{
				{1, "title1", "description1"},
				{2, "title2", "description2"},
				{3, "title3", "description3"},
			},
		},
		{
			name: "Error",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM lists l INNER JOIN users_lists ul ON (.+) WHERE (.+)").WithArgs(1).WillReturnError(errors.New("some error"))
			},
			userID:  1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := list.GetAll(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListPostgres_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	list := NewListPostgres(db)
	type args struct {
		userID int
		listID int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    bookshelf.List
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(1, "title", "description")
				mock.ExpectQuery("SELECT (.+) FROM lists l INNER JOIN users_lists ul ON (.+) WHERE (.+)").WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				userID: 1,
				listID: 1,
			},
			want: bookshelf.List{
				ID:          1,
				Title:       "title",
				Description: "description",
			},
			wantErr: false,
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})
				mock.ExpectQuery("SELECT (.+) FROM lists l INNER JOIN users_lists ul ON (.+) WHERE (.+)").WithArgs(1, 2).WillReturnRows(rows)
			},
			input: args{
				userID: 1,
				listID: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := list.GetByID(tt.input.userID, tt.input.listID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	list := NewListPostgres(db)
	type args struct {
		userID   int
		listID   int
		listItem bookshelf.List
		input    bookshelf.UpdateListInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.
					ExpectExec("UPDATE lists l SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("title2", "description2", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				userID: 1,
				listID: 1,
				listItem: bookshelf.List{
					ID:          1,
					Title:       "title",
					Description: "description",
				},
				input: bookshelf.UpdateListInput{
					Title:       stringPointer("title2"),
					Description: stringPointer("description2"),
				},
			},
		},
		{
			name: "Empty input",
			mock: func() {
				mock.
					ExpectExec("UPDATE lists l SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("title", "description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				userID: 1,
				listID: 1,
				listItem: bookshelf.List{
					ID:          1,
					Title:       "title",
					Description: "description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := list.Update(tt.input.userID, tt.input.listID, tt.input.listItem, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	list := NewListPostgres(db)
	type args struct {
		userID int
		listID int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.
					ExpectExec("DELETE FROM lists l USING users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				userID: 1,
				listID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := list.Delete(tt.input.userID, tt.input.listID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
