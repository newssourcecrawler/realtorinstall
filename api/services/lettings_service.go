package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type LettingsService struct {
	repo repos.LettingsRepo
}

func NewLettingsService(r repos.LettingsRepo) *LettingsService {
	return &LettingsService{repo: r}
}

func (s *LettingsService) CreateLetting(
	ctx context.Context,
	tenantID string,
	currentUser string,
	l models.Lettings,
) (int64, error) {
	// Required fields
	if l.PropertyID == 0 || l.TenantUserID == 0 || l.RentAmount <= 0 {
		return 0, errors.New("property, tenant user, and rent amount must be specified")
	}
	now := time.Now().UTC()
	l.TenantID = tenantID
	l.CreatedAt = now
	l.LastModified = now
	l.CreatedBy = currentUser
	l.ModifiedBy = currentUser
	l.Deleted = false

	return s.repo.Create(ctx, &l)
}

func (s *LettingsService) ListLettings(
	ctx context.Context,
	tenantID string,
) ([]models.Lettings, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Lettings, 0, len(rows))
	for _, rec := range rows {
		out = append(out, *rec)
	}
	return out, nil
}

func (s *LettingsService) UpdateLetting(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
	l models.Lettings,
) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	l.ID = id
	l.TenantID = tenantID
	l.CreatedAt = existing.CreatedAt
	l.CreatedBy = existing.CreatedBy
	l.Deleted = existing.Deleted
	l.ModifiedBy = currentUser
	l.LastModified = now

	return s.repo.Update(ctx, &l)
}

func (s *LettingsService) DeleteLetting(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
) error {
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
