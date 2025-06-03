package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

type PropertyService struct {
	repo        repos.PropertyRepo
	pricingRepo repos.LocationPricingRepo
}

// NewPropertyService constructs a new PropertyService.
func NewPropertyService(r repos.PropertyRepo, pr repos.LocationPricingRepo) *PropertyService {
	return &PropertyService{repo: r, pricingRepo: pr}
}

// CreateProperty accepts a models.Property whose date fields are strings in RFC3339 format.
// If the incoming model does not set those strings, we populate them here before saving.
func (s *PropertyService) CreateProperty(ctx context.Context, p models.Property) (int64, error) {
	if p.Address == "" || p.City == "" {
		return 0, errors.New("address and city cannot be empty")
	}

	// If ListingDate is empty, set it to now (RFC3339 string)
	if p.ListingDate == "" {
		p.ListingDate = time.Now().Format(time.RFC3339)
	}
	// Set created/modified timestamps to now
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	p.LastModified = now

	return s.repo.Create(ctx, &p)
}

// ListProperties returns all properties. It assumes the repo has stored date fields as RFC3339 strings.
func (s *PropertyService) ListProperties(ctx context.Context) ([]models.Property, error) {
	props, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Property
	for _, p := range props {
		// Each p.ListingDate, p.CreatedAt, p.LastModified is already a string
		out = append(out, *p)
	}
	return out, nil
}

// UpdateProperty allows editing of Address, City, or ZIP. It also updates LastModified.
func (s *PropertyService) UpdateProperty(ctx context.Context, p models.Property) error {
	if p.ID == 0 {
		return errors.New("invalid property ID")
	}
	// Update the LastModified field to now
	p.LastModified = time.Now().Format(time.RFC3339)
	return s.repo.Update(ctx, &p)
}
