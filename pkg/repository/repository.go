package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
)

type Book interface {
	Create(userId int, book models.Book) (int, error)
	GetAll() ([]models.Book, error)
	GetByUserId(userId int) ([]models.Book, error)
	Delete(bookId int) error
	Update(bookId int, input models.UpdateBookInput) error
}

type Comment interface {
	Create(userId int, book models.Comment) (int, error)
	GetAll() ([]models.Comment, error)
	GetByBookId(bookId int) ([]models.Comment, error)
	Delete(commentId int) error
	Update(commentId int, input models.UpdateCommentInput) error
}

type Repository struct {
	Book
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Book:    NewBookPostgres(db),
		Comment: NewCommentPostgres(db),
	}
}
