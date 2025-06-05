package services

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type ReportService struct {
	commissionRepo repos.CommissionRepo
	salesRepo      repos.SalesRepo
	pricingRepo    repos.LocationPricingRepo
	// … add more repos if needed
}

func NewReportService(
	cr repos.CommissionRepo,
	sr repos.SalesRepo,
	pr repos.LocationPricingRepo,
) *ReportService {
	return &ReportService{
		commissionRepo: cr,
		salesRepo:      sr,
		pricingRepo:    pr,
	}
}

// CommissionSummary holds beneficiary + total commission
type CommissionSummary struct {
	BeneficiaryID   int64   `json:"beneficiary_id"`
	TotalCommission float64 `json:"total_commission"`
}

func (s *ReportService) CommissionByBeneficiary(ctx context.Context, tenantID string) ([]CommissionSummary, error) {
	// Directly query the view_commission_by_beneficiary.
	// Assume your CommissionRepo exposes a RawQuery method (or just grab the *sql.DB).
	rows, err := s.commissionRepo.QueryCommissionByBeneficiary(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CommissionSummary
	for rows.Next() {
		var summary CommissionSummary
		if err := rows.Scan(&summary.BeneficiaryID, &summary.TotalCommission); err != nil {
			return nil, err
		}
		out = append(out, summary)
	}
	return out, nil
}
