package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type CommissionService struct {
	repo        repos.CommissionRepo
	saleRepo    repos.SaleRepo         // or repos.InstallmentPlanRepo if you treat sales as installments
	lettingRepo repos.LettingRepo      // if you have a separate letting table
	introRepo   repos.IntroductionRepo // if you have an introduction table
	userRepo    repos.UserRepo
}

func NewCommissionService(
	cr repos.CommissionRepo,
	sr repos.SaleRepo,
	lr repos.LettingRepo,
	ir repos.IntroductionRepo,
	ur repos.UserRepo,
) *CommissionService {
	return &CommissionService{
		repo:        cr,
		saleRepo:    sr,
		lettingRepo: lr,
		introRepo:   ir,
		userRepo:    ur,
	}
}

func (s *CommissionService) CreateCommission(
	ctx context.Context,
	tenantID string,
	currentUser string,
	comm models.Commission,
) (int64, error) {
	// Must have required fields
	if comm.TransactionType == "" || comm.TransactionID == 0 || comm.BeneficiaryID == 0 || comm.CommissionType == "" {
		return 0, errors.New("missing required commission fields")
	}
	// Validate transaction exists (optional)
	switch comm.TransactionType {
	case "sale":
		if _, err := s.saleRepo.GetByID(ctx, tenantID, comm.TransactionID); err != nil {
			return 0, errors.New("sale not found")
		}
	case "letting":
		if _, err := s.lettingRepo.GetByID(ctx, tenantID, comm.TransactionID); err != nil {
			return 0, errors.New("letting not found")
		}
	case "introduction":
		if _, err := s.introRepo.GetByID(ctx, tenantID, comm.TransactionID); err != nil {
			return 0, errors.New("introduction not found")
		}
	default:
		return 0, errors.New("invalid transaction type")
	}
	// Validate beneficiary exists
	if _, err := s.userRepo.GetByID(ctx, tenantID, comm.BeneficiaryID); err != nil {
		return 0, errors.New("beneficiary user not found")
	}

	// Calculate “calculated_amount” if percentage
	if comm.CommissionType == "percentage" {
		var txnValue float64
		switch comm.TransactionType {
		case "sale":
			sale, _ := s.saleRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = sale.TotalPrice
		case "letting":
			let, _ := s.lettingRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = let.RentalAmount
		case "introduction":
			intro, _ := s.introRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = intro.FeeAmount
		}
		comm.CalculatedAmount = txnValue * comm.RateOrAmount
	} else {
		// For “fixed” or “credit”, assume RateOrAmount already holds the exact number
		comm.CalculatedAmount = comm.RateOrAmount
	}

	now := time.Now().UTC()
	comm.TenantID = tenantID
	comm.CreatedAt = now
	comm.LastModified = now
	comm.CreatedBy = currentUser
	comm.ModifiedBy = currentUser
	comm.Deleted = false

	return s.repo.Create(ctx, &comm)
}

func (s *CommissionService) ListCommissions(
	ctx context.Context,
	tenantID string,
	filterType string, // optional; “sale” | “letting” | “introduction” | “” for all
	beneficiaryID int64, // optional; 0 = no filter
) ([]models.Commission, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Commission, 0, len(rows))
	for _, c := range rows {
		if c.Deleted {
			continue
		}
		if filterType != "" && c.TransactionType != filterType {
			continue
		}
		if beneficiaryID != 0 && c.BeneficiaryID != beneficiaryID {
			continue
		}
		out = append(out, *c)
	}
	return out, nil
}

// GET /reports/commissions/by‐beneficiary
func (h *ReportHandler) CommissionsByBeneficiary(c *gin.Context) {
	tenantID := c.GetString("currentTenant")
	data, err := h.svc.GetCommissionsByBeneficiary(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s *CommissionService) DeleteCommission(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
) error {
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
