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

func (s *PricingService) CreateLocationPricing(ctx context.Context, lp models.LocationPricing) (int64, error) {
	if lp.ZipCode == "" {
		return 0, errors.New("zip_code is required")
	}
	lp.CreatedAt = time.Now().UTC()
	lp.LastModified = lp.CreatedAt
	return s.repo.Create(ctx, &lp)
}

func (s *PricingService) ListLocationPricings(ctx context.Context) ([]models.LocationPricing, error) {
	lps, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.LocationPricing
	for _, lp := range lps {
		out = append(out, *lp)
	}
	return out, nil
}
