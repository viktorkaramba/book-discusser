package repository

import (
	"book-discusser/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, role) values ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	row, err := r.db.Query(query, email, password)
	defer row.Close()
	user = models.User{}
	for row.Next() {
		err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, errors.New("not found")
	}
	return &user, err
}

func (r *AuthPostgres) GetUserById(userId int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	row, err := r.db.Query(query, userId)
	defer row.Close()
	user = models.User{}
	for row.Next() {
		err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, errors.New("not found")
	}
	return &user, err
}

func (r *AuthPostgres) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1", usersTable)
	row, err := r.db.Query(query, email)
	defer row.Close()
	user = models.User{}
	for row.Next() {
		err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, errors.New("not found")
	}
	return &user, err
}
