package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCommentPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type args struct {
		userId  int
		bookId  int
		comment models.Comment
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
				bookId: 1,
				comment: models.Comment{
					Message: "Good Book!!!",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs(args.comment.Message).WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO books_comments").
					WithArgs(args.userId, args.bookId, id).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "2nd Insert Error",
			args: args{
				userId: 1,
				bookId: 1,
				comment: models.Comment{
					Message: "Good Book!!!",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs(args.comment.Message).WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO books_comments").WithArgs(args.bookId, id).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty Fields",
			args: args{
				userId: 1,
				bookId: 1,
				comment: models.Comment{
					Message: "",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("empty field"))
				mock.ExpectQuery("INSERT INTO comments").
					WithArgs(args.comment.Message).WillReturnRows(rows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)
			got, err := r.Create(testCase.args.userId, testCase.args.bookId, testCase.args.comment)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestCommentPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []models.Comment
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message"}).
					AddRow(1, "message1").
					AddRow(2, "message2").
					AddRow(3, "message3")

				mock.ExpectQuery("SELECT (.+) FROM comments c").WillReturnRows(rows)
			},
			want: []models.Comment{
				{1, "message1"},
				{2, "message2"},
				{3, "message3"},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message"})

				mock.ExpectQuery("SELECT (.+) FROM comments c").WillReturnRows(rows)
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

func TestCommentPostgres_GetByBookId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type args struct {
		bookId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.UsersComments
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message", "name"}).
					AddRow(1, "message1", "user1").
					AddRow(2, "message2", "user2")

				mock.ExpectQuery("SELECT (.+) FROM comments c JOIN books_comments bc on c.id = bc.comment_id " +
					"JOIN users u ON u.id = bc.User_id WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				bookId: 1,
			},
			want: []models.UsersComments{
				{1, "message1", "user1"},
				{2, "message2", "user2"},
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message", "name"})

				mock.ExpectQuery("SELECT (.+) FROM comments c JOIN books_comments bc on c.id = bc.comment_id " +
					"JOIN users u ON u.id = bc.User_id WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
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

			got, err := r.GetByBookId(tt.input.bookId)
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

func TestCommentPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type args struct {
		commentId int
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
				mock.ExpectExec("DELETE FROM comments WHERE (.+)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				commentId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM comments WHERE (.+)").
					WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				commentId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.commentId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCommentPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCommentPostgres(db)

	type args struct {
		commentId int
		input     models.UpdateCommentInput
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
				mock.ExpectExec("UPDATE comments c SET (.+) WHERE (.+)").
					WithArgs("new message", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				commentId: 1,
				input: models.UpdateCommentInput{
					Message: stringPointer("new message"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE comments c SET WHERE (.+)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				commentId: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.commentId, tt.input.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
