package models

import "errors"

type Book struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name" binding:"required"`
	Author    string `json:"author" db:"author"  binding:"required"`
	ImageBook string `json:"imageBook" db:"imagebook" binding:"required"`
}

type Comment struct {
	ID      int    `json:"id" db:"id"`
	Message string `json:"message" db:"message"  binding:"required"`
}

type UsersBook struct {
	ID     int
	UserId int
	BookId int
}

type BooksComment struct {
	ID        int
	BookId    int
	CommentId int
}

type UserCreateBook struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	ImageBook string `json:"imageBook"`
}

type UpdateBookInput struct {
	Name      *string `json:"name"`
	Author    *string `json:"author"`
	ImageBook *string `json:"imageBook"`
}

func (i UpdateBookInput) Validate() error {
	if i.Name == nil && i.Author == nil && i.ImageBook == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateCommentInput struct {
	Message *string `json:"message"`
}

func (i UpdateCommentInput) Validate() error {
	if i.Message == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
