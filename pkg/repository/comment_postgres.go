package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type CommentPostgres struct {
	db *sql.DB
}

func NewCommentPostgres(db *sql.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (c *CommentPostgres) Create(bookId int, comment models.Comment) (int, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return 0, err
	}

	var commentId int
	createCommentQuery := fmt.Sprintf("INSERT INTO %s (message) values ($1) RETURNING id", commentsTable)

	row := tx.QueryRow(createCommentQuery, comment.Message)
	err = row.Scan(&commentId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createBooksCommentQuery := fmt.Sprintf("INSERT INTO %s (book_Id, comment_Id) values ($1, $2)", booksCommentsTable)
	_, err = tx.Exec(createBooksCommentQuery, bookId, commentId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return commentId, tx.Commit()
}

func (c *CommentPostgres) GetAll() ([]models.Comment, error) {
	var comments []models.Comment
	query := fmt.Sprintf(`SELECT c.id, c.message FROM %s c`, commentsTable)
	rowsRs, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rowsRs.Close()
	for rowsRs.Next() {
		comment := models.Comment{}
		err := rowsRs.Scan(&comment.ID, &comment.Message)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (c *CommentPostgres) GetByBookId(bookId int) ([]models.Comment, error) {
	var comments []models.Comment
	query := fmt.Sprintf(`SELECT c.id, c.message FROM %s c
								INNER JOIN %s bc on c.id = bc.comment_id WHERE bc.book_id = $1`,
		commentsTable, booksCommentsTable)
	row, err := c.db.Query(query, bookId)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		comment := models.Comment{}
		err := row.Scan(&comment.ID, &comment.Message)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if len(comments) == 0 {
		return nil, errors.New("not found")
	}
	return comments, err
}

func (c *CommentPostgres) Delete(commentId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", commentsTable)
	_, err := c.db.Exec(query, commentId)
	return err
}

func (c *CommentPostgres) Update(commentId int, input models.UpdateCommentInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Message != nil {
		setValues = append(setValues, fmt.Sprintf("message=$%d", argId))
		args = append(args, *input.Message)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s c SET %s WHERE c.id = $1`, commentsTable, setQuery)
	args = append(args, commentId)
	_, err := c.db.Exec(query, args...)
	return err
}
