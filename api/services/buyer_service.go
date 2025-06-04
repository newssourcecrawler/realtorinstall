// services/buyer_service.go
package services

import (
	"context"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type BuyerService struct {
	repo repos.BuyerRepo
}

func NewBuyerService(r repos.BuyerRepo) *BuyerService {
	return &BuyerService{repo: r}
}

func (s *BuyerService) CreateBuyer(ctx context.Context, tenantID, currentUser string, b models.Buyer) (int64, error) {
	if b.FirstName == "" || b.LastName == "" || b.Email == "" {
		return 0, repos.ErrNameEmailNotFound
	}
	now := time.Now().UTC()
	b.TenantID = tenantID
	b.CreatedAt = now
	b.LastModified = now
	b.CreatedBy = currentUser
	b.ModifiedBy = currentUser
	b.Deleted = false
	return s.repo.Create(ctx, &b)
}

func (s *BuyerService) ListBuyers(ctx context.Context, tenantID string) ([]models.Buyer, error) {
	bs, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Buyer, 0, len(bs))
	for _, b := range bs {
		if b.Deleted {
			continue
		}
		out = append(out, *b)
	}
	return out, nil
}

func (s *BuyerService) UpdateBuyer(ctx context.Context, tenantID, currentUser string, id int64, b models.Buyer) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	b.TenantID = tenantID
	b.ID = id
	b.ModifiedBy = currentUser
	b.LastModified = now
	return s.repo.Update(ctx, &b)
}

func (s *BuyerService) DeleteBuyer(ctx context.Context, tenantID, currentUser string, id int64) error {
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
