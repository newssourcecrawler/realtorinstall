package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type CommissionService struct {
	repo        repos.CommissionRepo
	saleRepo    repos.SalesRepo
	lettingRepo repos.LettingsRepo
	introRepo   repos.IntroductionsRepo
	userRepo    repos.UserRepo
}

func NewCommissionService(
	cr repos.CommissionRepo,
	sr repos.SalesRepo,
	lr repos.LettingsRepo,
	ir repos.IntroductionsRepo,
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
	// Required fields
	if comm.TransactionType == "" ||
		comm.TransactionID == 0 ||
		comm.BeneficiaryID == 0 ||
		comm.CommissionType == "" {
		return 0, errors.New("missing required commission fields")
	}

	// Validate referenced transaction
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

	// Validate beneficiary
	if _, err := s.userRepo.GetByID(ctx, tenantID, comm.BeneficiaryID); err != nil {
		return 0, errors.New("beneficiary not found")
	}

	// Compute calculated_amount if percentage
	if comm.CommissionType == "percentage" {
		var txnValue float64
		switch comm.TransactionType {
		case "sale":
			sale, _ := s.saleRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = sale.SalePrice
		case "letting":
			let, _ := s.lettingRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = let.RentAmount
		case "introduction":
			intro, _ := s.introRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = intro.AgreedFee
		}
		comm.CalculatedAmount = txnValue * comm.RateOrAmount
	} else {
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
	filterType string, // “sale”|“letting”|“introduction” or “” for all
	beneficiaryID int64, // 0 for no filter
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

func (s *CommissionService) UpdateCommission(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
	comm models.Commission,
) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}

	// Preserve fields that shouldn’t change
	comm.ID = id
	comm.TenantID = tenantID
	comm.CreatedAt = existing.CreatedAt
	comm.CreatedBy = existing.CreatedBy
	comm.Deleted = existing.Deleted

	// Recompute calculated_amount if percentage
	if comm.CommissionType == "percentage" {
		var txnValue float64
		switch comm.TransactionType {
		case "sale":
			sale, _ := s.saleRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = sale.SalePrice
		case "letting":
			let, _ := s.lettingRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = let.RentAmount
		case "introduction":
			intro, _ := s.introRepo.GetByID(ctx, tenantID, comm.TransactionID)
			txnValue = intro.AgreedFee
		}
		comm.CalculatedAmount = txnValue * comm.RateOrAmount
	} else {
		comm.CalculatedAmount = comm.RateOrAmount
	}

	now := time.Now().UTC()
	comm.ModifiedBy = currentUser
	comm.LastModified = now

	return s.repo.Update(ctx, &comm)
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
