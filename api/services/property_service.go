package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("not found")

type PropertyService struct {
	repo        repos.PropertyRepo
	pricingRepo repos.LocationPricingRepo
}

// NewPropertyService constructs a new PropertyService.
func NewPropertyService(r repos.PropertyRepo, pr repos.LocationPricingRepo) *PropertyService {
	return &PropertyService{repo: r, pricingRepo: pr}
}

// CreateProperty creates a new property, setting listing & timestamp fields if empty.
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

// ListProperties returns all properties from the repo.
func (s *PropertyService) ListProperties(ctx context.Context) ([]models.Property, error) {
	props, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Property
	for _, p := range props {
		out = append(out, *p)
	}
	return out, nil
}

// UpdateProperty edits an existing property. Returns ErrNotFound if the repo signals no match.
func (s *PropertyService) UpdateProperty(ctx context.Context, id string, p models.Property) error {
	// Convert id (string) to int64 inside the repo layer; assume the repo.Update returns ErrNotFound when not found.
	if p.ID == 0 {
		return IDNotFound
	}
	// Refresh LastModified
	p.LastModified = time.Now().Format(time.RFC3339)

	err := s.repo.Update(ctx, &p)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}

// DeleteProperty removes a property by ID. Returns ErrNotFound if not found.
func (s *PropertyService) DeleteProperty(ctx context.Context, id string) error {
	// The repo.Delete method should handle converting id→int64 and return ErrNotFound if needed.
	err := s.repo.Delete(ctx, id)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}
