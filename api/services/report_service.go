package services

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type ReportService struct {
	commissionRepo       repos.CommissionRepo
	installmentsplanRepo repos.InstallmentPlanRepo
	salesRepo            repos.SalesRepo
	lettingsRepo         repos.LettingsRepo
	propertyRepo         repos.PropertyRepo
	// (drop salesRepo and pricingRepo for now; add back later if you bake in more reports)
}

func NewReportService(cr repos.CommissionRepo) *ReportService {
	return &ReportService{commissionRepo: cr}
}

func (s *ReportService) CommissionByBeneficiary(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.commissionRepo.SummarizeByBeneficiary(ctx, tenantID)
}

func (s *ReportService) OutstandingInstallmentsByPlan(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.installmentsplanRepo.SummarizeByPlan(ctx, tenantID)
}

func (s *ReportService) MonthlySalesVolume(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.salesRepo.SummarizeByMonth(ctx, tenantID)
}

func (s *ReportService) ActiveLettingsRentRoll(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.lettingsRepo.SummarizeRentRoll(ctx, tenantID)
}

func (s *ReportService) TopPropertiesByPaymentVolume(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.propertyRepo.SummarizeTopProperties(ctx, tenantID)
}
