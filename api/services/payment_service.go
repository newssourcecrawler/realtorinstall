package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

type PaymentService struct {
	repo               repos.PaymentRepo
	installmentService *InstallmentService
}

func NewPaymentService(r repos.PaymentRepo, ir *InstallmentService) *PaymentService {
	return &PaymentService{repo: r, installmentService: ir}
}

func (s *PaymentService) CreatePayment(ctx context.Context, p models.Payment) (int64, error) {
	if p.InstallmentID == 0 {
		return 0, errors.New("installment_id is required")
	}
	p.CreatedAt = time.Now().UTC()
	p.LastModified = p.CreatedAt
	return s.repo.Create(ctx, &p)
}

func (s *PaymentService) ListPayments(ctx context.Context) ([]models.Payment, error) {
	ps, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Payment
	for _, p := range ps {
		out = append(out, *p)
	}
	return out, nil
}
