// services/installment_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type InstallmentService struct {
	repo        repos.InstallmentRepo
	paymentRepo repos.PaymentRepo
}

func NewInstallmentService(r repos.InstallmentRepo, pr repos.PaymentRepo) *InstallmentService {
	return &InstallmentService{repo: r, paymentRepo: pr}
}

func (s *InstallmentService) CreateInstallment(ctx context.Context, tenantID, currentUser string, inst models.Installment) (int64, error) {
	if inst.PlanID == 0 {
		return 0, errors.New("plan_id is required")
	}
	now := time.Now().UTC()
	inst.TenantID = tenantID
	inst.CreatedAt = now
	inst.LastModified = now
	inst.CreatedBy = currentUser
	inst.ModifiedBy = currentUser
	inst.Deleted = false
	return s.repo.Create(ctx, &inst)
}

func (s *InstallmentService) ListInstallments(ctx context.Context, tenantID string) ([]models.Installment, error) {
	is, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Installment, 0, len(is))
	for _, i := range is {
		out = append(out, *i)
	}
	return out, nil
}

func (s *InstallmentService) ListByPlan(ctx context.Context, tenantID string, planID int64) ([]models.Installment, error) {
	is, err := s.repo.ListByPlan(ctx, tenantID, planID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Installment, 0, len(is))
	for _, i := range is {
		out = append(out, *i)
	}
	return out, nil
}

func (s *InstallmentService) UpdateInstallment(ctx context.Context, tenantID, currentUser string, id int64, inst models.Installment) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}
	now := time.Now().UTC()
	inst.TenantID = tenantID
	inst.ID = id
	inst.ModifiedBy = currentUser
	inst.LastModified = now
	return s.repo.Update(ctx, &inst)
}

func (s *InstallmentService) DeleteInstallment(ctx context.Context, tenantID, currentUser string, id int64) error {
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
