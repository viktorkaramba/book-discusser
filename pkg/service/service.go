package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
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

type Service struct {
	Authorization
	Book
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Book:          NewBookService(repos.Book),
		Comment:       NewCommentService(repos.Comment),
	}
}
