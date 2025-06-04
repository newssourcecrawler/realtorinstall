package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type PlanService struct {
	repo        repos.InstallmentPlanRepo
	installRepo repos.InstallmentRepo
}

func NewPlanService(r repos.InstallmentPlanRepo, ir repos.InstallmentRepo) *PlanService {
	return &PlanService{repo: r, installRepo: ir}
}

func (s *PlanService) CreatePlan(ctx context.Context, p models.InstallmentPlan) (int64, error) {
	if p.PropertyID == 0 || p.BuyerID == 0 {
		return 0, errors.New("property and buyer must be specified")
	}
	p.CreatedAt = time.Now().UTC()
	p.LastModified = p.CreatedAt
	return s.repo.Create(ctx, &p)
}

func (s *PlanService) ListPlans(ctx context.Context) ([]models.InstallmentPlan, error) {
	ps, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.InstallmentPlan
	for _, p := range ps {
		out = append(out, *p)
	}
	return out, nil
}
