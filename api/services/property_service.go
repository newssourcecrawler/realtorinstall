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
	userRepo    repos.UserRepo // to look up username if needed
	pricingRepo repos.LocationPricingRepo
}

// NewPropertyService constructs a new PropertyService.
func NewPropertyService(r repos.PropertyRepo, u repos.UserRepo, pr repos.LocationPricingRepo) *PropertyService {
	return &PropertyService{repo: r, userRepo: u, pricingRepo: pr}
}

// CreateProperty creates a new property, setting listing & timestamp fields if empty.
func (s *PropertyService) CreateProperty(ctx context.Context, tenantID, currentUser string, p models.Property) (int64, error) {
	if p.Address == "" || p.City == "" || p.ZIP == "" {
		return 0, repos.ErrAddrNotFound
	}

	// If ListingDate is empty, set it to now (RFC3339 string)
	if p.ListingDate.IsZero() {
		//p.ListingDate = time.Now().Format(time.RFC3339)
		p.ListingDate = time.Now().UTC()
	}
	// Set created/modified timestamps to now
	//now := time.Now().Format(time.RFC3339)
	p.TenantID = tenantID
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now
	p.CreatedBy = currentUser
	p.ModifiedBy = currentUser
	p.Deleted = false

	return s.repo.Create(ctx, &p)
}

// ListProperties returns all properties from the repo.
func (s *PropertyService) ListProperties(ctx context.Context, tenantID string) ([]models.Property, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	var out []models.Property
	for _, p := range rows {
		out = append(out, *p)
	}
	return out, nil
}

// UpdateProperty edits an existing property. Returns ErrNotFound if the repo signals no match.
func (s *PropertyService) UpdateProperty(ctx context.Context, tenantID, currentUser string, id int64, p models.Property) error {
	// Convert id (string) to int64 inside the repo layer; assume the repo.Update returns ErrNotFound when not found.
	if p.ID == 0 {
		return repos.ErrIDNotFound
	}
	// Refresh LastModified
	//p.LastModified = time.Now().Format(time.RFC3339)
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

// DeleteProperty performs a soft‐delete (marks Deleted=true). Returns ErrNotFound if missing.
func (s *PropertyService) DeleteProperty(ctx context.Context, tenantID, currentUser string, id int64) error {
	// The repo.Delete method will look up the existing row and set deleted=1
	return s.repo.Delete(ctx, tenantID, id)
}
