package services

import (
	"context"
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
func (s *PropertyService) CreateProperty(ctx context.Context, p models.Property, createdByID int64) (int64, error) {
	if p.Address == "" || p.City == "" {
		return 0, repos.ErrAddrNotFound
	}

	// If ListingDate is empty, set it to now (RFC3339 string)
	if p.ListingDate.IsZero() {
		//p.ListingDate = time.Now().Format(time.RFC3339)
		p.ListingDate = time.Now().UTC()
	}
	// Set created/modified timestamps to now
	//now := time.Now().Format(time.RFC3339)
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now

	// Look up username of createdByID if you want to store the name rather than ID
	user, err := s.userRepo.GetByID(ctx, createdByID) // userRepo injected earlier
	if err != nil {
		return 0, err
	}
	p.CreatedBy = user.Username
	p.ModifiedBy = user.Username
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
func (s *PropertyService) UpdateProperty(ctx context.Context, id int64, p models.Property, modifiedByID int64) error {
	// Convert id (string) to int64 inside the repo layer; assume the repo.Update returns ErrNotFound when not found.
	if p.ID == 0 {
		return repos.ErrIDNotFound
	}
	// Refresh LastModified
	//p.LastModified = time.Now().Format(time.RFC3339)
	p.LastModified = time.Now().UTC()

	//err := s.repo.Update(ctx, &p)
	user, err := s.userRepo.GetByID(ctx, modifiedByID)
	if err != nil {
		return err
	}
	p.ModifiedBy = user.Username

	err = s.repo.Update(ctx, &p)
	if err == repos.ErrNotFound {
		return repos.ErrNotFound
	}
	return err
}

// DeleteProperty performs a soft‐delete (marks Deleted=true). Returns ErrNotFound if missing.
func (s *PropertyService) DeleteProperty(ctx context.Context, id int64, p models.Property, modifiedByID int64) error {
	// 1) Load the existing record (even if it was previously soft‐deleted)
	existing, err := s.repo.GetByID(ctx, id)
	if err == repos.ErrNotFound {
		return repos.ErrNotFound
	} else if err != nil {
		return err
	}

	// 2) Mark it as deleted
	existing.Deleted = true
	existing.LastModified = time.Now().UTC()
	existing.ModifiedBy = user.Username

	err = s.repo.Update(ctx, existing)
	if err == repos.ErrNotFound {
		return repos.ErrNotFound
	}
	return err
}
