package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
)

type BookService struct {
	bookRepo    repository.Book
	commentRepo repository.Comment
}

func NewBookService(bookRepo repository.Book, commentRepo repository.Comment) *BookService {
	return &BookService{bookRepo: bookRepo, commentRepo: commentRepo}
}

func (b *BookService) Create(userId int, userBookComment models.UserCreateBook) (int, error) {
	book := models.Book{Name: userBookComment.Name, Author: userBookComment.Author,
		Description: userBookComment.Description, ImageBook: userBookComment.ImageBook}
	comment := models.Comment{Message: userBookComment.Message}
	id, err := b.bookRepo.Create(userId, book)
	if err != nil {
		return -1, err
	}
	_, err = b.commentRepo.Create(userId, id, comment)
	return id, err
}

func (b *BookService) GetAll() ([]models.Book, error) {
	return b.bookRepo.GetAll()
}

func (b *BookService) GetByUserId(userId int) ([]models.Book, error) {
	return b.bookRepo.GetByUserId(userId)
}

func (b *BookService) Delete(bookId int) error {
	return b.bookRepo.Delete(bookId)
}

func (b *BookService) Update(bookId int, input models.UpdateBookInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return b.bookRepo.Update(bookId, input)
}
