// services/auth_service.go
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

func (s *AuthService) Register(ctx context.Context, tenantID, currentUser string, u models.User, rawPassword string) (int64, error) {
	if u.UserName == "" || u.Role == "" || u.FirstName == "" || u.LastName == "" {
		return 0, errors.New("username, role, first and last name required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.PasswordHash = string(hash)
	now := time.Now().UTC()
	u.TenantID = tenantID
	u.CreatedAt = now
	u.LastModified = now
	u.CreatedBy = currentUser
	u.ModifiedBy = currentUser
	u.Deleted = false
	return s.repo.Create(ctx, &u)
}

func (s *AuthService) Authenticate(ctx context.Context, tenantID, username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(ctx, tenantID, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
