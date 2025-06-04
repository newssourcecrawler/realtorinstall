package services

import (
	"context"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) HealthCheck(ctx context.Context) (string, error) {
	return "OK", nil
}
