package services

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type ReportService struct {
	commissionRepo      repos.CommissionRepo
	installmentPlanRepo repos.InstallmentPlanRepo
	salesRepo           repos.SalesRepo
	lettingsRepo        repos.LettingsRepo
	propertyRepo        repos.PropertyRepo
}

func NewReportService(
	cr repos.CommissionRepo,
	ipr repos.InstallmentPlanRepo,
	sr repos.SalesRepo,
	lr repos.LettingsRepo,
	pr repos.PropertyRepo,
) *ReportService {
	return &ReportService{
		commissionRepo:      cr,
		installmentPlanRepo: ipr,
		salesRepo:           sr,
		lettingsRepo:        lr,
		propertyRepo:        pr,
	}
}

func (s *ReportService) TotalCommissionByBeneficiary(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	return s.commissionRepo.TotalCommissionByBeneficiary(ctx, tenantID)
}

func (s *ReportService) OutstandingInstallmentsByPlan(ctx context.Context, tenantID string) ([]models.PlanSummary, error) {
	return s.installmentPlanRepo.SummarizeByPlan(ctx, tenantID)
}

func (s *ReportService) MonthlySalesVolume(ctx context.Context, tenantID string) ([]models.MonthSales, error) {
	return s.salesRepo.SummarizeByMonth(ctx, tenantID)
}

func (s *ReportService) ActiveLettingsRentRoll(ctx context.Context, tenantID string) ([]models.RentRoll, error) {
	return s.lettingsRepo.SummarizeRentRoll(ctx, tenantID)
}

func (s *ReportService) TopPropertiesByPaymentVolume(ctx context.Context, tenantID string) ([]models.PropertyPaymentVolume, error) {
	return s.propertyRepo.SummarizeTopProperties(ctx, tenantID)
}
