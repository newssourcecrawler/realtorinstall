package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService struct {
	repo repos.UserRepo
}

func NewAuthService(r repos.UserRepo) *AuthService {
	return &AuthService{repo: r}
}

// Register user: hash password, create record
func (s *AuthService) Register(ctx context.Context, u models.User, rawPassword string) (int64, error) {
	if u.UserName == "" || u.PasswordHash == "" || u.Role == "" || u.FirstName == "" || u.LastName == "" {
		return 0, errors.New("username, password, first and last name and role required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.PasswordHash = string(hash)
	now := time.Now().UTC()
	u.CreatedAt = now
	u.LastModified = now
	u.Deleted = false
	return s.repo.Create(ctx, &u)
}

// Authenticate look up username, compare password
func (s *AuthService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
