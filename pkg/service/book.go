package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

func (b *BookService) Create(userId int, book models.Book) (int, error) {
	return b.repo.Create(userId, book)
}

func (b *BookService) GetAll() ([]models.Book, error) {
	return b.repo.GetAll()
}

func (b *BookService) GetByUserId(userId int) ([]models.Book, error) {
	return b.repo.GetByUserId(userId)
}

func (b *BookService) Delete(bookId int) error {
	return b.repo.Delete(bookId)
}

func (b *BookService) Update(bookId int, input models.UpdateBookInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return b.repo.Update(bookId, input)
}
