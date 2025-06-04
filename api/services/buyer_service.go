package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

// ErrNotFound is already defined in this package; reuse it.
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
	b.CreatedAt = time.Now().Format(time.RFC3339)
	b.LastModified = b.CreatedAt
	return s.repo.Create(ctx, &b)
}

func (s *BuyerService) ListBuyers(ctx context.Context) ([]models.Buyer, error) {
	bs, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	var out []models.Buyer
	for _, b := range bs {
		out = append(out, *b)
	}
	return out, nil
}

// UpdateBuyer edits an existing buyer. Returns ErrNotFound if the repo signals no match.
func (s *BuyerService) UpdateBuyer(ctx context.Context, id string, b models.Buyer) error {
	// We ignore any ID in 'b' and rely on the repo.Update to use 'id' string.
	b.LastModified = time.Now().Format(time.RFC3339)
	err := s.repo.Update(ctx, id, &b)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}

// DeleteBuyer removes a buyer by ID. Returns ErrNotFound if not found.
func (s *BuyerService) DeleteBuyer(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == repos.ErrNotFound {
		return ErrNotFound
	}
	return err
}
