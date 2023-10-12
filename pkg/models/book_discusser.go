package models

type Book struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Author string `json:"author" db:"author"`
}

type Comment struct {
	ID      int    `json:"id" db:"id"`
	Message string `json:"message" db:"message"`
}

type UsersBook struct {
	ID     int
	UserId int
	BookId int
}

type BooksComment struct {
	ID          int
	BookId      int
	UserComment int
}
