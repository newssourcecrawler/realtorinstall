package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("not found")

type BuyerService struct {
	repo repos.BuyerRepo
}

func NewBuyerService(r repos.BuyerRepo) *BuyerService {
	return &BuyerService{repo: r}
}

func (s *BuyerService) CreateBuyer(ctx context.Context, b models.Buyer) (int64, error) {
	if b.Name == "" || b.Email == "" {
		return 0, errors.New("name and email are required")
	}
	b.Deleted = false
	now := time.Now().Format(time.RFC3339)
	b.CreatedAt = now
	b.LastModified = now
	return s.repo.Create(ctx, &b)
}

func (s *BuyerService) ListBuyers(ctx context.Context) ([]models.Buyer, error) {
	bs, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Buyer
	for _, b := range bs {
		if b.Deleted {
			continue
		}
		out = append(out, *b)
	}
	return out, nil
}

// UpdateBuyer edits an existing buyer. Returns ErrNotFound if the repo signals no match.
func (s *BuyerService) UpdateBuyer(ctx context.Context, id string, b models.Buyer) error {
	b.LastModified = time.Now().Format(time.RFC3339)
	err := s.repo.Update(ctx, id, &b)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}

// DeleteBuyer performs a soft‐delete (marks 'Deleted=true') instead of hard deletion.
func (s *BuyerService) DeleteBuyer(ctx context.Context, id string) error {
	b, err := s.repo.GetByID(ctx, id)
	if err == repos.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	b.Deleted = true
	b.LastModified = time.Now().Format(time.RFC3339)
	err = s.repo.Update(ctx, id, b)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}
