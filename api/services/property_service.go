package services

import (
	"context"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type PropertyService struct {
	repo        repos.PropertyRepo
	userRepo    repos.UserRepo
	pricingRepo repos.LocationPricingRepo
}

func NewPropertyService(r repos.PropertyRepo, u repos.UserRepo, pr repos.LocationPricingRepo) *PropertyService {
	return &PropertyService{repo: r, userRepo: u, pricingRepo: pr}
}

func (s *PropertyService) CreateProperty(ctx context.Context, tenantID, currentUser string, p models.Property) (int64, error) {
	if p.Address == "" || p.City == "" || p.ZIP == "" {
		return 0, repos.ErrAddrNotFound
	}
	if p.ListingDate.IsZero() {
		p.ListingDate = time.Now().UTC()
	}
	now := time.Now().UTC()
	p.TenantID = tenantID
	p.CreatedAt = now
	p.LastModified = now
	p.CreatedBy = currentUser
	p.ModifiedBy = currentUser
	p.Deleted = false
	return s.repo.Create(ctx, &p)
}

func (s *PropertyService) ListProperties(ctx context.Context, tenantID string) ([]models.Property, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Property, 0, len(rows))
	for _, p := range rows {
		out = append(out, *p)
	}
	return out, nil
}

func (s *PropertyService) UpdateProperty(ctx context.Context, tenantID, currentUser string, id int64, p models.Property) error {
	if p.ID == 0 {
		return repos.ErrIDNotFound
	}
	p.TenantID = tenantID
	p.ID = id
	p.ModifiedBy = currentUser
	p.LastModified = time.Now().UTC()
	err := s.repo.Update(ctx, &p)
	if err == repos.ErrNotFound {
		return repos.ErrNotFound
	}
	return err
}

func (s *PropertyService) DeleteProperty(ctx context.Context, tenantID, currentUser string, id int64) error {
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	p.Deleted = true
	p.ModifiedBy = currentUser
	p.LastModified = time.Now().UTC()
	return s.repo.Update(ctx, p)
}
