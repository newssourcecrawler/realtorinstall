package services

import (
	"context"
	"errors"
	"time"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
	"github.com/newssourcecrawler/realtorinstall/internal/repos"
)

type BuyerService struct {
	repo repos.BuyerRepo
}

func NewBuyerService(r repos.BuyerRepo) *BuyerService {
	return &BuyerService{repo: r}
}

func (s *BuyerService) CreateBuyer(ctx context.Context, b models.Buyer) (int64, error) {
	if b.Name == "" {
		return 0, errors.New("name is required")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	b.CreatedAt = now
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
