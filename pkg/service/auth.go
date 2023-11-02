package service

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/repository"
	"crypto/sha1"
	"fmt"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
)

type AuthService struct {
	repoAuth repository.Authorization
}

func NewAuthService(repoAuth repository.Authorization) *AuthService {
	return &AuthService{repoAuth: repoAuth}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repoAuth.CreateUser(user)
}

func (s *AuthService) GetUser(email, password string) (*models.User, error) {
	return s.repoAuth.GetUser(email, generatePasswordHash(password))
}

func (s *AuthService) GetUserById(userId int) (*models.User, error) {
	return s.repoAuth.GetUserById(userId)
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.repoAuth.GetUserByEmail(email)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
