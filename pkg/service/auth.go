package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
	"book-discusser/pkg/sessions"
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
)

type AuthService struct {
	repoAuth    repository.Authorization
	repoSession repository.Session
}

func NewAuthService(repoAuth repository.Authorization, repoSession repository.Session) *AuthService {
	return &AuthService{repoAuth: repoAuth, repoSession: repoSession}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repoAuth.CreateUser(user)
}

func (s *AuthService) GenerateSessionToken(userId int, email, password string) (*sessions.Session, error) {
	_, err := s.repoAuth.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return nil, err
	}
	sessionToken := uuid.NewString()
	newSession := sessions.Session{
		ID:     sessionToken,
		UserId: userId,
		Email:  email,
		Expiry: time.Now().Add(sessions.ExpireTime),
	}
	_, err = s.repoSession.CreateSession(newSession)
	if err != nil {
		return nil, err
	}
	return &newSession, nil
}

func (s *AuthService) GetSession(sessionId string) (*sessions.Session, error) {
	return s.repoSession.GetSession(sessionId)
}

func (s *AuthService) DeleteSession(sessionId string) error {
	return s.repoSession.Delete(sessionId)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
