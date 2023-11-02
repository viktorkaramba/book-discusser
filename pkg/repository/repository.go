package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (*models.User, error)
	GetUserById(userId int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type Book interface {
	Create(userId int, book models.Book) (int, error)
	GetAll() ([]models.Book, error)
	GetByUserId(userId int) ([]models.Book, error)
	Delete(bookId int) error
	Update(bookId int, input models.UpdateBookInput) error
}

type Comment interface {
	Create(userId, bookId int, book models.Comment) (int, error)
	GetAll() ([]models.Comment, error)
	GetByBookId(bookId int) ([]models.UsersComments, error)
	GetByUserId(userId int) ([]models.Comment, error)
	Delete(commentId int) error
	Update(commentId int, input models.UpdateCommentInput) error
}

type Repository struct {
	Authorization
	Book
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Book:          NewBookPostgres(db),
		Comment:       NewCommentPostgres(db),
	}
}
