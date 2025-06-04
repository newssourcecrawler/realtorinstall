package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// ErrNotFound is returned when a record does not exist.
//var ErrNotFound = errors.New("not found")

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
	if p.ListingDate == "" {
		p.ListingDate = time.Now().Format(time.RFC3339)
	}
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	p.LastModified = now
	p.Deleted = false
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
		if p.Deleted {
			continue
		}
		out = append(out, *p)
	}
	return out, nil
}

// UpdateProperty edits an existing property. Returns ErrNotFound if not found.
func (s *PropertyService) UpdateProperty(ctx context.Context, id string, p models.Property) error {
	p.LastModified = time.Now().Format(time.RFC3339)
	// We expect repo.Update to return repos.ErrNotFound if no row exists.
	err := s.repo.Update(ctx, &p)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}

// DeleteProperty performs a soft‐delete (marks 'Deleted=true') instead of hard‐deletion.
func (s *PropertyService) DeleteProperty(ctx context.Context, id string) error {
	prop, err := s.repo.GetByID(ctx, id)
	if err == repos.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	prop.Deleted = true
	prop.LastModified = time.Now().Format(time.RFC3339)
	err = s.repo.Update(ctx, prop)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}
