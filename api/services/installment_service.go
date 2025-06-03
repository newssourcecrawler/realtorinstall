package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

type InstallmentService struct {
	repo        repos.InstallmentRepo
	paymentRepo repos.PaymentRepo
}

func NewInstallmentService(r repos.InstallmentRepo, pr repos.PaymentRepo) *InstallmentService {
	return &InstallmentService{repo: r, paymentRepo: pr}
}

func (s *InstallmentService) CreateInstallment(ctx context.Context, inst models.Installment) (int64, error) {
	if inst.PlanID == 0 {
		return 0, errors.New("plan_id is required")
	}
	inst.CreatedAt = time.Now().UTC()
	inst.LastModified = inst.CreatedAt
	return s.repo.Create(ctx, &inst)
}

func (s *InstallmentService) ListInstallments(ctx context.Context) ([]models.Installment, error) {
	is, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Installment
	for _, i := range is {
		out = append(out, *i)
	}
	return out, nil
}
