package repository

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/sessions"
	"database/sql"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (*models.User, error)
}

type Session interface {
	CreateSession(session sessions.Session) (string, error)
	GetSession(sessionId string) (*sessions.Session, error)
	Delete(sessionId string) error
}

type Book interface {
	Create(userId int, book models.Book) (int, error)
	GetAll() ([]models.Book, error)
	GetByUserId(userId int) ([]models.Book, error)
	Delete(bookId int) error
	Update(bookId int, input models.UpdateBookInput) error
}

type Comment interface {
	Create(bookId int, book models.Comment) (int, error)
	GetAll() ([]models.Comment, error)
	GetByBookId(bookId int) ([]models.Comment, error)
	Delete(commentId int) error
	Update(commentId int, input models.UpdateCommentInput) error
}

type Repository struct {
	Authorization
	Book
	Comment
	Session
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Book:          NewBookPostgres(db),
		Comment:       NewCommentPostgres(db),
		Session:       NewSessionPostgres(db),
	}
}
