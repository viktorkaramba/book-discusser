package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (c *CommentService) Create(userId, bookId int, comment models.Comment) (int, error) {
	return c.repo.Create(userId, bookId, comment)
}

func (c *CommentService) GetAll() ([]models.Comment, error) {
	return c.repo.GetAll()
}

func (c *CommentService) GetByBookId(bookId int) ([]models.UsersComments, error) {
	return c.repo.GetByBookId(bookId)
}

func (c *CommentService) Delete(commentId int) error {
	return c.repo.Delete(commentId)
}

func (c *CommentService) Update(commentId int, input models.UpdateCommentInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return c.repo.Update(commentId, input)
}
