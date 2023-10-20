package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (b *BookPostgres) Create(userId int, book models.Book) (int, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	var bookId int
	createBookQuery := fmt.Sprintf("INSERT INTO %s (name, Author, imagebook) values ($1, $2, $3) RETURNING id", booksTable)

	row := tx.QueryRow(createBookQuery, book.Name, book.Author, book.ImageBook)
	err = row.Scan(&bookId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserBooksQuery := fmt.Sprintf("INSERT INTO %s (user_Id, book_Id) values ($1, $2)", usersBooksTable)
	_, err = tx.Exec(createUserBooksQuery, userId, bookId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return bookId, tx.Commit()
}

func (b *BookPostgres) GetAll() ([]models.Book, error) {
	var books []models.Book
	query := fmt.Sprintf(`SELECT b.id, b.name, b.author, b.imagebook FROM %s b`, booksTable)
	rowsRs, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rowsRs.Close()
	for rowsRs.Next() {
		book := models.Book{}
		err := rowsRs.Scan(&book.ID, &book.Name, &book.Author, &book.ImageBook)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (b *BookPostgres) GetByUserId(userId int) ([]models.Book, error) {
	var books []models.Book
	query := fmt.Sprintf(`SELECT b.id, b.name, b.author, b.imageBook FROM %s b
								INNER JOIN %s ub on b.id = ub.book_id WHERE ub.user_id = $1`,
		booksTable, usersBooksTable)
	row, err := b.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		book := models.Book{}
		err := row.Scan(&book.ID, &book.Name, &book.Author, &book.ImageBook)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if len(books) == 0 {
		return nil, errors.New("not found")
	}
	return books, err
}

func (b *BookPostgres) Delete(bookId int) error {
	query := fmt.Sprintf("DELETE FROM %s b WHERE b.id=$1", booksTable)
	_, err := b.db.Exec(query, bookId)

	return err
}

func (b *BookPostgres) Update(bookId int, input models.UpdateBookInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, *input.Author)
		argId++
	}

	if input.ImageBook != nil {
		setValues = append(setValues, fmt.Sprintf("imageBook=$%d", argId))
		args = append(args, *input.Author)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s b SET %s WHERE b.id = $1`,
		booksTable, setQuery)
	args = append(args, bookId)

	_, err := b.db.Exec(query, args...)
	return err
}
