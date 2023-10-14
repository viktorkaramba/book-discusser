package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBookPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		userId int
		book   models.Book
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				userId: 1,
				book: models.Book{
					Name:   "Verity",
					Author: "Colleen Hoover",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO books").
					WithArgs(args.book.Name, args.book.Author).WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO users_books").
					WithArgs(args.userId, id).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "2nd Insert Error",
			args: args{
				userId: 1,
				book: models.Book{
					Name:   "Verity",
					Author: "Colleen Hoover",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO books").
					WithArgs(args.book.Name, args.book.Author).WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO users_books").WithArgs(args.userId, id).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty Fields",
			args: args{
				userId: 1,
				book: models.Book{
					Name:   "",
					Author: "Colleen Hoover",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("empty field"))
				mock.ExpectQuery("INSERT INTO books").
					WithArgs(args.book.Name, args.book.Author).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)
			got, err := r.Create(testCase.args.userId, testCase.args.book)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestBookPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []models.Book
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "author"}).
					AddRow(1, "name1", "author1").
					AddRow(2, "name2", "author2").
					AddRow(3, "name3", "author3")

				mock.ExpectQuery("SELECT (.+) FROM books b").WillReturnRows(rows)
			},
			want: []models.Book{
				{1, "name1", "author1"},
				{2, "name2", "author2"},
				{3, "name3", "author3"},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "author"})

				mock.ExpectQuery("SELECT (.+) FROM books b").WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll()
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

func TestBookPostgres_GetByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.Book
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "author"}).
					AddRow(1, "name1", "author1").
					AddRow(2, "name2", "author2")

				mock.ExpectQuery("SELECT (.+) FROM books b INNER JOIN users_books ub on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: []models.Book{
				{1, "name1", "author1"},
				{2, "name2", "author2"},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "author"})

				mock.ExpectQuery("SELECT (.+) FROM books b INNER JOIN users_books ub on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetByUserId(tt.input.userId)
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

func TestBookPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		bookId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM books b WHERE (.+)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				bookId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM books b WHERE (.+)").
					WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				bookId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.bookId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		bookId int
		input  models.UpdateBookInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE books b SET (.+) WHERE (.+)").
					WithArgs("new name", "new author", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				bookId: 1,
				input: models.UpdateBookInput{
					Name:   stringPointer("new name"),
					Author: stringPointer("new author"),
				},
			},
		},
		{
			name: "OK_WithoutAuthor",
			mock: func() {
				mock.ExpectExec("UPDATE books b SET (.+) WHERE (.+)").
					WithArgs("new name", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				bookId: 1,
				input: models.UpdateBookInput{
					Name: stringPointer("new name"),
				},
			},
		},
		{
			name: "OK_WithoutName",
			mock: func() {
				mock.ExpectExec("UPDATE books b SET (.+) WHERE (.+)").
					WithArgs("new author", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				bookId: 1,
				input: models.UpdateBookInput{
					Author: stringPointer("new author"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE books b SET WHERE (.+)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				bookId: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.bookId, tt.input.input)
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
