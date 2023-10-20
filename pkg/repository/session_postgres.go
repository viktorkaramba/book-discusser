package repository

import (
	"book-discusser/pkg/sessions"
	"database/sql"
	"errors"
	"fmt"
)

type SessionPostgres struct {
	db *sql.DB
}

func NewSessionPostgres(db *sql.DB) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (r *SessionPostgres) CreateSession(session sessions.Session) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (id, userId, email, expiry) values ($1, $2, $3, $4) RETURNING id", sessionsTable)
	row := r.db.QueryRow(query, session.ID, session.UserId, session.Email, session.Expiry)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *SessionPostgres) GetSession(sessionId string) (*sessions.Session, error) {
	var session sessions.Session
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", sessionsTable)
	row, err := r.db.Query(query, sessionId)
	defer row.Close()
	session = sessions.Session{}
	for row.Next() {
		err = row.Scan(&session.ID, &session.Email, &session.Expiry)
		if err != nil {
			return nil, err
		}
	}
	if session.ID == "" {
		return nil, errors.New("not found")
	}
	return &session, err
}

func (r *SessionPostgres) Delete(sessionId string) error {
	query := fmt.Sprintf("DELETE FROM %s s WHERE s.id=$1", sessionsTable)
	_, err := r.db.Exec(query, sessionId)
	return err
}
