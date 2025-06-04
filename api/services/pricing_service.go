// services/pricing_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type PricingService struct {
	repo repos.LocationPricingRepo
}

func NewPricingService(r repos.LocationPricingRepo) *PricingService {
	return &PricingService{repo: r}
}

func (s *PricingService) CreateLocationPricing(ctx context.Context, tenantID, currentUser string, lp models.LocationPricing) (int64, error) {
	if lp.ZipCode == "" {
		return 0, errors.New("zip_code is required")
	}
	now := time.Now().UTC()
	lp.TenantID = tenantID
	lp.CreatedAt = now
	lp.LastModified = now
	lp.CreatedBy = currentUser
	lp.ModifiedBy = currentUser
	lp.Deleted = false
	return s.repo.Create(ctx, &lp)
}

func (s *PricingService) ListLocationPricings(ctx context.Context, tenantID string) ([]models.LocationPricing, error) {
	lps, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.LocationPricing, 0, len(lps))
	for _, lp := range lps {
		out = append(out, *lp)
	}
	return out, nil
}

func (s *PricingService) UpdateLocationPricing(ctx context.Context, tenantID, currentUser string, id int64, lp models.LocationPricing) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	lp.TenantID = tenantID
	lp.ID = id
	lp.ModifiedBy = currentUser
	lp.LastModified = now
	return s.repo.Update(ctx, &lp)
}

func (s *PricingService) DeleteLocationPricing(ctx context.Context, tenantID, currentUser string, id int64) error {
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
