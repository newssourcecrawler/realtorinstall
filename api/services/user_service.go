// services/User_service.go
package services

import (
	"context"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type UserService struct {
	repo repos.UserRepo
}

func NewUserService(r repos.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, tenantID, currentUser string, b models.User) (int64, error) {
	if b.FirstName == "" || b.LastName == "" || b.Email == "" {
		return 0, repos.ErrNameEmailNotFound
	}
	now := time.Now().UTC()
	b.TenantID = tenantID
	b.CreatedAt = now
	b.LastModified = now
	b.CreatedBy = currentUser
	b.ModifiedBy = currentUser
	b.Deleted = false
	return s.repo.Create(ctx, &b)
}

func (s *UserService) ListUsers(ctx context.Context, tenantID string) ([]models.User, error) {
	bs, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.User, 0, len(bs))
	for _, b := range bs {
		if b.Deleted {
			continue
		}
		out = append(out, *b)
	}
	return out, nil
}

func (s *UserService) UpdateUser(ctx context.Context, tenantID, currentUser string, id int64, b models.User) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	b.TenantID = tenantID
	b.ID = id
	b.ModifiedBy = currentUser
	b.LastModified = now
	return s.repo.Update(ctx, &b)
}

func (s *UserService) DeleteUser(ctx context.Context, tenantID, currentUser string, id int64) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	existing.Deleted = true
	existing.ModifiedBy = currentUser
	existing.LastModified = time.Now().UTC()
	return s.repo.Update(ctx, existing)
}
