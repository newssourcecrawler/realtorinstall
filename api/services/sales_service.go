// services/Sales_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type SalesService struct {
	repo        repos.InstallmentSalesRepo
	installRepo repos.InstallmentRepo
}

func NewSalesService(r repos.InstallmentSalesRepo, ir repos.InstallmentRepo) *SalesService {
	return &SalesService{repo: r, installRepo: ir}
}

func (s *SalesService) CreateSales(ctx context.Context, tenantID, currentUser string, p models.InstallmentSales) (int64, error) {
	if p.PropertyID == 0 || p.BuyerID == 0 {
		return 0, errors.New("property and buyer must be specified")
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

func (s *SalesService) ListSaless(ctx context.Context, tenantID string) ([]models.InstallmentSales, error) {
	ps, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.InstallmentSales, 0, len(ps))
	for _, p := range ps {
		out = append(out, *p)
	}
	return out, nil
}

func (s *SalesService) UpdateSales(ctx context.Context, tenantID, currentUser string, id int64, p models.InstallmentSales) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	p.TenantID = tenantID
	p.ID = id
	p.ModifiedBy = currentUser
	p.LastModified = now
	return s.repo.Update(ctx, &p)
}

func (s *SalesService) DeleteSales(ctx context.Context, tenantID, currentUser string, id int64) error {
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
