// services/payment_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type PaymentService struct {
	repo               repos.PaymentRepo
	installmentService *InstallmentService
}

func NewPaymentService(r repos.PaymentRepo, ir *InstallmentService) *PaymentService {
	return &PaymentService{repo: r, installmentService: ir}
}

func (s *PaymentService) CreatePayment(ctx context.Context, tenantID, currentUser string, p models.Payment) (int64, error) {
	if p.InstallmentID == 0 {
		return 0, errors.New("installment_id is required")
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

func (s *PaymentService) ListPayments(ctx context.Context, tenantID string) ([]models.Payment, error) {
	ps, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Payment, 0, len(ps))
	for _, p := range ps {
		out = append(out, *p)
	}
	return out, nil
}

func (s *PaymentService) UpdatePayment(ctx context.Context, tenantID, currentUser string, id int64, p models.Payment) error {
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

func (s *PaymentService) DeletePayment(ctx context.Context, tenantID, currentUser string, id int64) error {
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
