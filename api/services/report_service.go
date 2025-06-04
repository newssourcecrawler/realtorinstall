package services

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type ReportService struct {
	installRepo repos.InstallmentRepo
}

func NewReportService(ir repos.InstallmentRepo) *ReportService {
	return &ReportService{installRepo: ir}
}

// Stub: in the future, you might implement methods like OverdueInstallments, UpcomingDue, etc.
func (s *ReportService) HealthCheck(ctx context.Context) (string, error) {
	return "OK", nil
}
