package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

type IntroductionsService struct {
	repo        repos.IntroductionsRepo
	saleRepo    repos.SalesRepo
	lettingRepo repos.LettingsRepo
	userRepo    repos.UserRepo
}

func NewIntroductionsService(
	r repos.IntroductionsRepo,
	sr repos.SalesRepo,
	lr repos.LettingsRepo,
	ur repos.UserRepo,
) *IntroductionsService {
	return &IntroductionsService{
		repo:        r,
		saleRepo:    sr,
		lettingRepo: lr,
		userRepo:    ur,
	}
}

func (s *IntroductionsService) CreateIntroduction(
	ctx context.Context,
	tenantID string,
	currentUser string,
	intro models.Introductions,
) (int64, error) {
	// Required fields
	if intro.IntroducerID == 0 ||
		intro.IntroducedParty == "" ||
		intro.PropertyID == 0 {
		return 0, errors.New("introducer, introduced party, and property must be specified")
	}

	// If this introduction is already linked to a TransactionID, verify existence
	if intro.TransactionID != 0 {
		switch intro.TransactionType {
		case "sale":
			if _, err := s.saleRepo.GetByID(ctx, tenantID, intro.TransactionID); err != nil {
				return 0, errors.New("sale not found")
			}
		case "letting":
			if _, err := s.lettingRepo.GetByID(ctx, tenantID, intro.TransactionID); err != nil {
				return 0, errors.New("letting not found")
			}
		default:
			return 0, errors.New("invalid transaction type")
		}
	}

	// Validate introducer exists
	if _, err := s.userRepo.GetByID(ctx, tenantID, intro.IntroducerID); err != nil {
		return 0, errors.New("introducer user not found")
	}

	now := time.Now().UTC()
	intro.TenantID = tenantID
	intro.CreatedAt = now
	intro.LastModified = now
	intro.CreatedBy = currentUser
	intro.ModifiedBy = currentUser
	intro.Deleted = false

	return s.repo.Create(ctx, &intro)
}

func (s *IntroductionsService) ListIntroductions(
	ctx context.Context,
	tenantID string,
) ([]models.Introductions, error) {
	rows, err := s.repo.ListAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]models.Introductions, 0, len(rows))
	for _, rec := range rows {
		out = append(out, *rec)
	}
	return out, nil
}

func (s *IntroductionsService) UpdateIntroduction(
	ctx context.Context,
	tenantID string,
	currentUser string,
	id int64,
	intro models.Introductions,
) error {
	existing, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return repos.ErrNotFound
	}

	// Preserve unchangeable fields
	intro.ID = id
	intro.TenantID = tenantID
	intro.CreatedAt = existing.CreatedAt
	intro.CreatedBy = existing.CreatedBy
	intro.Deleted = existing.Deleted

	// If TransactionID present, verify again
	if intro.TransactionID != 0 {
		switch intro.TransactionType {
		case "sale":
			if _, err := s.saleRepo.GetByID(ctx, tenantID, intro.TransactionID); err != nil {
				return errors.New("sale not found")
			}
		case "letting":
			if _, err := s.lettingRepo.GetByID(ctx, tenantID, intro.TransactionID); err != nil {
				return errors.New("letting not found")
			}
		default:
			return errors.New("invalid transaction type")
		}
	}

	now := time.Now().UTC()
	intro.ModifiedBy = currentUser
	intro.LastModified = now

	return s.repo.Update(ctx, &intro)
}

func (s *IntroductionsService) DeleteIntroduction(
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
