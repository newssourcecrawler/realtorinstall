package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type SalesService struct {
	repo repos.SalesRepo
}

func NewSalesService(r repos.SalesRepo) *SalesService {
	return &SalesService{repo: r}
}

func (s *SalesService) CreateSale(
	ctx context.Context,
	tenantID string,
	currentUser string,
	sale models.Sales,
) (int64, error) {
	// Required fields
	if sale.PropertyID == 0 || sale.BuyerID == 0 || sale.SalePrice <= 0 || sale.SaleType == "" {
		return 0, errors.New("property, buyer, sale price, and sale type must be specified")
	}
	now := time.Now().UTC()
	sale.TenantID = tenantID
	sale.CreatedAt = now
	sale.LastModified = now
	sale.CreatedBy = currentUser
	sale.ModifiedBy = currentUser
	sale.Deleted = false

	return s.repo.Create(ctx, &sale)
}

func (s *SalesService) ListSales(
	ctx context.Context,
	tenantID string,
) ([]models.Sales, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Sales, 0, len(rows))
	for _, rec := range rows {
		out = append(out, *rec)
	}
	return out, nil
}

func (s *SalesService) UpdateSale(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
	sale models.Sales,
) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}

	// Preserve unchangeable fields
	sale.ID = id
	sale.TenantID = tenantID
	sale.CreatedAt = existing.CreatedAt
	sale.CreatedBy = existing.CreatedBy
	sale.Deleted = existing.Deleted

	now := time.Now().UTC()
	sale.ModifiedBy = currentUser
	sale.LastModified = now

	return s.repo.Update(ctx, &sale)
}

func (s *SalesService) DeleteSale(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
) error {
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
