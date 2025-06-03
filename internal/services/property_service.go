// internal/services/property_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

// PropertyService wraps PropertyRepo + pricingRepo
type PropertyService struct {
	repo        repos.PropertyRepo
	pricingRepo repos.LocationPricingRepo
}

// NewPropertyService constructs one
func NewPropertyService(r repos.PropertyRepo, pr repos.LocationPricingRepo) *PropertyService {
	return &PropertyService{repo: r, pricingRepo: pr}
}

// CreateProperty creates and returns the new property ID
func (s *PropertyService) CreateProperty(ctx context.Context, p models.Property) (int64, error) {
	if p.Address == "" || p.City == "" {
		return 0, errors.New("address and city are required")
	}
	p.ListingDate = time.Now().UTC()
	return s.repo.Create(ctx, &p)
}

// ListProperties returns all properties (flattened)
func (s *PropertyService) ListProperties(ctx context.Context) ([]models.Property, error) {
	ps, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Property
	for _, p := range ps {
		out = append(out, *p)
	}
	return out, nil
}
